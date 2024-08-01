package controllers

import (
	"bos_personal_ai/config"

	"github.com/gofiber/fiber/v2"
)

type AskAIInput struct {
	Prompt string `json:"prompt" xml:"prompt" form:"prompt"`
}

type AskRAGInput struct {
	ModelProvider string `json:"model_provider" xml:"model_provider" form:"model_provider"`
	Question      string `json:"question" xml:"question" form:"question"`
}

type AskAIController struct {
	app *config.AppConfig
}

func NewAskAIController(app *config.AppConfig) *AskAIController {
	return &AskAIController{
		app: app,
	}
}

func (ctrl *AskAIController) AskRAG(c *fiber.Ctx) error {
	input := new(AskRAGInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_error": true,
			"message":  "error when parsing payload: " + err.Error(),
		})
	}

	aiAnswer, err := ctrl.app.RagService.AskQuestion(input.ModelProvider, input.Question)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_error": true,
			"message":  "error when ask to AI: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"is_error": false,
		"data":     aiAnswer,
		"success":  true,
		"message":  "successfully created",
	})
}
