package infra_services

import "bos_personal_ai/thirdparties"

type AIChatServiceInterface interface {
	Prompt(serviceProvider string, systemPrompt string, userPrompt string, temp float32, maxToken int) (string, error)
}

type AIChatServiceImpl struct {
	openAIChat      thirdparties.AIChatInterface
	huggingfaceChat thirdparties.AIChatInterface
}

func NewAIChatService(
	openAIChat thirdparties.AIChatInterface,
	huggingfaceChat thirdparties.AIChatInterface,
) AIChatServiceInterface {
	return &AIChatServiceImpl{
		openAIChat:      openAIChat,
		huggingfaceChat: huggingfaceChat,
	}
}

func (srv *AIChatServiceImpl) Prompt(serviceProvider string, systemPrompt string, userPrompt string, temp float32, maxToken int) (string, error) {
	if serviceProvider == "huggingface" {
		return srv.huggingfaceChat.Prompt(systemPrompt, userPrompt, temp, maxToken)
	}
	return srv.openAIChat.Prompt(systemPrompt, userPrompt, temp, maxToken)
}
