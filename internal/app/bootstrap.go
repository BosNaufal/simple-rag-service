package bootstrap

import (
	"bos_personal_ai/config"
	"bos_personal_ai/repositories"
	app_services "bos_personal_ai/services/apps"
	infra_services "bos_personal_ai/services/infra"
	"bos_personal_ai/thirdparties"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Bootstrap() *config.AppConfig {
	// change this on production
	dsn := "host=localhost user=postgres password=root dbname=personal_ai port=5431 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	openAIEmbedding := thirdparties.NewEmbeddingOpenAIEmbedding()
	openAIChatService := thirdparties.NewOpenAIChatThirdParty()

	searchCacheRepo := repositories.NewSearchCacheRepository(db)
	searchCacheService := infra_services.NewSearchCacheService(searchCacheRepo)

	knowledgeRepo := repositories.NewKnowledgeRepository(db)
	knowledgeService := infra_services.NewKnowledgeService(knowledgeRepo)

	embeddedKnowledgeService := app_services.NewEmbeddedKnowledgeService(openAIEmbedding, knowledgeService, searchCacheService)
	ragService := app_services.NewRAG(openAIEmbedding, embeddedKnowledgeService)

	appConfig := config.AppConfig{
		ThirdParties: config.ThirdParties{
			Embedding: openAIEmbedding,
			AIChat:    openAIChatService,
		},
		KnowledgeServices:        knowledgeService,
		RagService:               ragService,
		EmbeddedKnowledgeService: embeddedKnowledgeService,
	}

	return &appConfig
}
