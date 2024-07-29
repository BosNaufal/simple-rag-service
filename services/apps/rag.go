package app_services

import (
	"bos_personal_ai/thirdparties"
)

type RAGInterface interface {
	AskQuestion(question string) (string, error)
}

type RAGImpl struct {
	embeddingThirdparty      thirdparties.EmbeddingThirdPartyInterface
	embeddedKnowledgeService EmbeddedKnowledgeServiceInterface
}

func NewRAG(
	embeddingThirdparty thirdparties.EmbeddingThirdPartyInterface,
	embeddedKnowledgeService EmbeddedKnowledgeServiceInterface,
) *RAGImpl {
	return &RAGImpl{
		embeddingThirdparty:      embeddingThirdparty,
		embeddedKnowledgeService: embeddedKnowledgeService,
	}
}

func (srv *RAGImpl) AskQuestion(question string) (string, error) {
	return "", nil
}
