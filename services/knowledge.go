package services

import (
	"bos_personal_ai/domain"
	"bos_personal_ai/repositories"
	"bos_personal_ai/thirdparties"

	"gorm.io/gorm"
)

type KnowledgeServicesInterface interface {
	Find(searchQuery string) ([]repositories.KnowledgeEntity, error)
	Add(title string, content string) (repositories.KnowledgeEntity, error)
	Update(input repositories.KnowledgeEntity) (repositories.KnowledgeEntity, error)
	Delete(id uint) error
}

type KnowledgeServicesImpl struct {
	embeddingService  thirdparties.EmbeddingThirdPartyInterface
	knowledgeDomain   domain.KnowledgeDomainInterface
	searchCacheDomain domain.SearchCacheDomainInterface
}

func NewKnowledgeService(
	embeddingServices thirdparties.EmbeddingThirdPartyInterface,
	knowledgeDomain domain.KnowledgeDomainInterface,
	searchCacheDomain domain.SearchCacheDomainInterface,
) *KnowledgeServicesImpl {
	return &KnowledgeServicesImpl{
		embeddingService:  embeddingServices,
		knowledgeDomain:   knowledgeDomain,
		searchCacheDomain: searchCacheDomain,
	}
}

func (srv *KnowledgeServicesImpl) Find(searchQuery string) ([]repositories.KnowledgeEntity, error) {
	searchCache, err := srv.searchCacheDomain.Find(searchQuery)
	if err != nil && err != gorm.ErrRecordNotFound {
		return []repositories.KnowledgeEntity{}, err
	}

	embeddingString := searchCache.Embedding
	if err == gorm.ErrRecordNotFound {
		embeddingString, err = srv.embeddingService.GetEmbeddingFromString(searchQuery)
		if err != nil {
			return []repositories.KnowledgeEntity{}, err
		}

		go func() {
			srv.searchCacheDomain.Add(searchQuery, embeddingString)
		}()
	}

	return srv.knowledgeDomain.Find(embeddingString)
}

func (srv *KnowledgeServicesImpl) getEmbeddingFromKnowledgeInput(title string, content string) (string, error) {
	concatContent := title + "\n\n" + content
	embeddingString, err := srv.embeddingService.GetEmbeddingFromString(concatContent)
	if err != nil {
		return embeddingString, err
	}
	return embeddingString, nil
}

func (srv *KnowledgeServicesImpl) Add(title string, content string) (repositories.KnowledgeEntity, error) {
	embeddingString, err := srv.getEmbeddingFromKnowledgeInput(title, content)
	if err != nil {
		return repositories.KnowledgeEntity{}, err
	}
	return srv.knowledgeDomain.Add(title, content, embeddingString)
}

func (srv *KnowledgeServicesImpl) Update(input repositories.KnowledgeEntity) (repositories.KnowledgeEntity, error) {
	newEmbeddingString, err := srv.getEmbeddingFromKnowledgeInput(input.Title, input.Content)
	if err != nil {
		return repositories.KnowledgeEntity{}, err
	}
	return srv.knowledgeDomain.Update(input, newEmbeddingString)
}

func (srv *KnowledgeServicesImpl) Delete(id uint) error {
	return srv.knowledgeDomain.Delete(id)
}
