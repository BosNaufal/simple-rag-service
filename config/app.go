package config

import (
	app_services "bos_personal_ai/services/apps"
	infra_services "bos_personal_ai/services/infra"
	thirdparties "bos_personal_ai/thirdparties"
)

type ThirdParties struct {
	Embedding thirdparties.EmbeddingThirdPartyInterface
}

type AppConfig struct {
	ThirdParties             ThirdParties
	KnowledgeServices        infra_services.KnowledgeServicesInterface
	RagService               app_services.RAGInterface
	EmbeddedKnowledgeService app_services.EmbeddedKnowledgeServiceInterface
}
