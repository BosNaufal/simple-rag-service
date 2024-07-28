package config

import (
	"bos_personal_ai/services"
	thirdparties "bos_personal_ai/thirdparties"
)

type ThirdParties struct {
	Embedding thirdparties.EmbeddingThirdPartyInterface
}

type AppConfig struct {
	ThirdParties      ThirdParties
	KnowledgeServices services.KnowledgeServicesInterface
}
