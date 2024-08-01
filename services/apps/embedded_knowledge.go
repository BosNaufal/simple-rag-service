package app_services

import (
	"bos_personal_ai/repositories"
	infra_services "bos_personal_ai/services/infra"
	"bos_personal_ai/thirdparties"

	"gorm.io/gorm"
)

type EmbeddedKnowledgeServiceInterface interface {
	RetriveKnowledgeBySearchQuery(searchQuery string) ([]repositories.KnowledgeEntity, error)
	AddNewKnowledgeWithEmbedding(title string, content string) (repositories.KnowledgeEntity, error)
	UpdateKnowledgeWithEmbedding(input repositories.KnowledgeEntity) (repositories.KnowledgeEntity, error)
}

type EmbeddedKnowledgeServiceImpl struct {
	embeddingThirdparty thirdparties.EmbeddingThirdPartyInterface
	knowledgeService    infra_services.KnowledgeServicesInterface
	searchCacheService  infra_services.SearchCachesInterface
}

func NewEmbeddedKnowledgeService(
	embeddingThirdparty thirdparties.EmbeddingThirdPartyInterface,
	knowledgeService infra_services.KnowledgeServicesInterface,
	searchCacheService infra_services.SearchCachesInterface,
) EmbeddedKnowledgeServiceInterface {
	return &EmbeddedKnowledgeServiceImpl{
		embeddingThirdparty: embeddingThirdparty,
		knowledgeService:    knowledgeService,
		searchCacheService:  searchCacheService,
	}
}

func (srv *EmbeddedKnowledgeServiceImpl) RetriveKnowledgeBySearchQuery(searchQuery string) ([]repositories.KnowledgeEntity, error) {
	searchCache, err := srv.searchCacheService.FindEmbeddingInCache(searchQuery)
	if err != nil && err != gorm.ErrRecordNotFound {
		return []repositories.KnowledgeEntity{}, err
	}

	embeddingString := searchCache.Embedding
	if err == gorm.ErrRecordNotFound {
		embeddingString, err = srv.embeddingThirdparty.GetEmbeddingFromString(searchQuery)
		if err != nil {
			return []repositories.KnowledgeEntity{}, err
		}

		go func() {
			srv.searchCacheService.SaveQueryInCache(searchQuery, embeddingString)
		}()
	}

	return srv.knowledgeService.FindByEmbedding(embeddingString)
}

func (srv *EmbeddedKnowledgeServiceImpl) getEmbeddingFromKnowledgeInput(title string, content string) (string, error) {
	concatContent := title + "\n\n" + content
	embeddingString, err := srv.embeddingThirdparty.GetEmbeddingFromString(concatContent)
	if err != nil {
		return embeddingString, err
	}
	return embeddingString, nil
}

func (srv *EmbeddedKnowledgeServiceImpl) AddNewKnowledgeWithEmbedding(title string, content string) (repositories.KnowledgeEntity, error) {
	embeddingString, err := srv.getEmbeddingFromKnowledgeInput(title, content)
	if err != nil {
		return repositories.KnowledgeEntity{}, err
	}
	return srv.knowledgeService.AddNewKnowledge(title, content, embeddingString)
}

func (srv *EmbeddedKnowledgeServiceImpl) UpdateKnowledgeWithEmbedding(input repositories.KnowledgeEntity) (repositories.KnowledgeEntity, error) {
	newEmbeddingString, err := srv.getEmbeddingFromKnowledgeInput(input.Title, input.Content)
	if err != nil {
		return repositories.KnowledgeEntity{}, err
	}

	return srv.knowledgeService.UpdateKnowledge(input, newEmbeddingString)
}
