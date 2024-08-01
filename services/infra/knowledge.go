package infra_services

import (
	"bos_personal_ai/repositories"
)

type KnowledgeServicesInterface interface {
	FindByEmbedding(embeddingString string) ([]repositories.KnowledgeEntity, error)
	AddNewKnowledge(title string, content string, embeddingString string) (repositories.KnowledgeEntity, error)
	UpdateKnowledge(input repositories.KnowledgeEntity, newEmbeddingString string) (repositories.KnowledgeEntity, error)
	DeleteKnowledge(id uint) error
}

type KnowledgeServicesImpl struct {
	knowledgeRepo repositories.KnowledgeRepositoryInterface
}

func NewKnowledgeService(
	knowledgeRepo repositories.KnowledgeRepositoryInterface,
) KnowledgeServicesInterface {
	return &KnowledgeServicesImpl{
		knowledgeRepo: knowledgeRepo,
	}
}

func (srv *KnowledgeServicesImpl) FindByEmbedding(embeddingString string) ([]repositories.KnowledgeEntity, error) {
	return srv.knowledgeRepo.Find(embeddingString, 5)
}

func (srv *KnowledgeServicesImpl) AddNewKnowledge(title string, content string, embeddingString string) (repositories.KnowledgeEntity, error) {
	return srv.knowledgeRepo.Add(repositories.KnowledgeEntity{
		Title:     title,
		Content:   content,
		Embedding: embeddingString,
	})
}

func (srv *KnowledgeServicesImpl) UpdateKnowledge(input repositories.KnowledgeEntity, newEmbeddingString string) (repositories.KnowledgeEntity, error) {
	input.Embedding = newEmbeddingString
	return srv.knowledgeRepo.Update(input)
}

func (srv *KnowledgeServicesImpl) DeleteKnowledge(id uint) error {
	return srv.knowledgeRepo.Delete(id)
}
