package controllers

import (
	"bos_personal_ai/config"
	"bos_personal_ai/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type KnowledgeInput struct {
	Title   string `json:"title" xml:"title" form:"title"`
	Content string `json:"content" xml:"content" form:"content"`
}

type KnowledgeController struct {
	app *config.AppConfig
}

func NewKnowledgeController(app *config.AppConfig) *KnowledgeController {
	return &KnowledgeController{
		app: app,
	}
}

func (ctrl *KnowledgeController) GetKnowledge(c *fiber.Ctx) error {
	searchQuery := c.Query("query")

	var knowledgeList []repositories.KnowledgeEntity
	var err error

	knowledgeList, err = ctrl.app.KnowledgeServices.Find(searchQuery)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_error": true,
			"message":  err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_error": false,
		"data":     knowledgeList,
		"success":  true,
		"message":  "successfully fetched",
	})
}

func (ctrl *KnowledgeController) CreateKnowledge(c *fiber.Ctx) error {
	input := new(KnowledgeInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_error": true,
			"message":  "error when parsing payload: " + err.Error(),
		})
	}

	result, err := ctrl.app.KnowledgeServices.Add(input.Title, input.Content)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_error": true,
			"message":  "error when creating data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"is_error": false,
		"data":     result,
		"success":  true,
		"message":  "successfully created",
	})
}

func (ctrl *KnowledgeController) DeleteKnowledge(c *fiber.Ctx) error {
	targetIdString := c.Params("id")
	targetIdUint, err := strconv.ParseUint(targetIdString, 10, 64)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed when parsing ID: " + err.Error(),
		})
	}

	err = ctrl.app.KnowledgeServices.Delete(uint(targetIdUint))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed when deleting the data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_error": false,
		"data":     nil,
		"success":  true,
		"message":  "successfully deleted",
	})
}

func (ctrl *KnowledgeController) UpdateKnowledge(c *fiber.Ctx) error {
	targetIdString := c.Params("id")
	targetIdUint, err := strconv.ParseUint(targetIdString, 10, 64)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed when parsing ID: " + err.Error(),
		})
	}

	input := new(repositories.KnowledgeEntity)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"is_error": true,
			"message":  "error when parsing payload: " + err.Error(),
		})
	}

	input.ID = uint(targetIdUint)
	updatedEntity, err := ctrl.app.KnowledgeServices.Update(*input)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"is_error": true,
			"message":  "failed when updating the data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"is_error": false,
		"data":     updatedEntity,
		"success":  true,
		"message":  "successfully deleted",
	})
}
