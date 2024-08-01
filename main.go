package main

import (
	controllers "bos_personal_ai/controllers/api"
	bootstrap "bos_personal_ai/internal/app"
	migration "bos_personal_ai/migrations"
	"log"
	"os"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
)

func main() {
	if len(os.Args) == 2 {
		firstArgument := os.Args[1]
		switch firstArgument {
		case "migrate-up":
			migration.MigrateUp()
			return
		case "migrate-down":
			migration.MigrateDown()
			return
		default:
			return
		}
	}

	app := fiber.New(fiber.Config{
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	appConfig := bootstrap.Bootstrap()

	api := app.Group("/api") // /api
	v1 := api.Group("/v1")   // /api/v1

	knowledgeController := controllers.NewKnowledgeController(appConfig)
	knowledge := v1.Group("/knowledge")
	knowledge.Get("/", knowledgeController.GetKnowledge)
	knowledge.Post("/", knowledgeController.CreateKnowledge)
	knowledge.Put("/:id", knowledgeController.UpdateKnowledge)
	knowledge.Delete("/:id", knowledgeController.DeleteKnowledge)

	askAIController := controllers.NewAskAIController(appConfig)
	ask := v1.Group("/ask")
	ask.Post("/rag", askAIController.AskRAG)

	log.Fatal(app.Listen(":3000"))
}
