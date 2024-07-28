package domain

import (
	"bos_personal_ai/repositories"
)

type KnowledgeDomainInterface interface {
	Find(searchQuery string) ([]repositories.KnowledgeEntity, error)
	Add(title string, content string, embeddingString string) (repositories.KnowledgeEntity, error)
	Update(input repositories.KnowledgeEntity, newEmbeddingString string) (repositories.KnowledgeEntity, error)
	Delete(id uint) error
}

type KnowledgeDomainImpl struct {
	knowledgeRepo repositories.KnowledgeRepositoryInterface
}

func NewKnowledgeDomain(kr *repositories.KnowledgeRepositoryImpl) *KnowledgeDomainImpl {
	return &KnowledgeDomainImpl{
		knowledgeRepo: kr,
	}
}

func (domain *KnowledgeDomainImpl) Find(embeddingString string) ([]repositories.KnowledgeEntity, error) {
	results, err := domain.knowledgeRepo.Find(embeddingString, 5)
	return results, err
}

func (domain *KnowledgeDomainImpl) Add(title string, content string, embeddingString string) (repositories.KnowledgeEntity, error) {
	knowledge := repositories.KnowledgeEntity{
		Title:     title,
		Content:   content,
		Embedding: embeddingString,
	}

	return domain.knowledgeRepo.Add(knowledge)
}

func (domain *KnowledgeDomainImpl) Update(input repositories.KnowledgeEntity, newEmbeddingString string) (repositories.KnowledgeEntity, error) {
	updatedData, err := domain.knowledgeRepo.Update(input)
	if err != nil {
		return repositories.KnowledgeEntity{}, err
	}

	return domain.knowledgeRepo.UpdateEmbedding(updatedData.ID, newEmbeddingString)
}

func (domain *KnowledgeDomainImpl) Delete(id uint) error {
	return domain.knowledgeRepo.Delete(id)
}
