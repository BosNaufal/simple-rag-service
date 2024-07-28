package bootstrap

import (
	"bos_personal_ai/config"
	"bos_personal_ai/domain"
	"bos_personal_ai/repositories"
	"bos_personal_ai/services"
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
	searchCacheRepo := repositories.NewSearchCacheRepository(db)
	searchCacheDomain := domain.NewSearchCacheDomain(searchCacheRepo)

	knowledgeRepo := repositories.NewKnowledgeRepository(db)
	knowledgeDomain := domain.NewKnowledgeDomain(
		knowledgeRepo,
	)
	knowledgeService := services.NewKnowledgeService(openAIEmbedding, knowledgeDomain, searchCacheDomain)

	appConfig := config.AppConfig{
		ThirdParties: config.ThirdParties{
			Embedding: openAIEmbedding,
		},
		KnowledgeServices: knowledgeService,
	}

	return &appConfig
}
