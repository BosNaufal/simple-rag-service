package thirdparties

import (
	app_constants "bos_personal_ai/env"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type AIChatInterface interface {
	Prompt(systemPrompt string, userPrompt string, temp float32, maxToken int) (string, error)
}

type OpenAIChatImpl struct {
}

type Choice struct {
	Index   int `json:"index"`
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           float64  `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

func NewOpenAIChatThirdParty() AIChatInterface {
	return &OpenAIChatImpl{}
}

func (e *OpenAIChatImpl) Prompt(systemPrompt string, userPrompt string, temp float32, maxToken int) (string, error) {
	var data CompletionResponse

	url := "https://api.openai.com/v1/chat/completions"

	postData := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]interface{}{
			{
				"role": "system",
				"content": []map[string]string{
					{
						"type": "text",
						"text": systemPrompt,
					},
				},
			},
			{
				"role": "user",
				"content": []map[string]string{
					{
						"type": "text",
						"text": userPrompt,
					},
				},
			},
		},
		"temperature":       temp,
		"max_tokens":        maxToken,
		"top_p":             1,
		"frequency_penalty": 0,
		"presence_penalty":  0,
	}

	jsonData, err := json.Marshal(postData)
	if err != nil {
		log.Fatalf("Failed to marshal JSON: %v", err)
		return "", err
	}

	// Create a new POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
		return "", err
	}

	// Set any required headers (if needed)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+app_constants.OPENAI_API_KEY)

	// Use the default HTTP client to send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send HTTP request: %v", err)
		return "", err
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
		return "", err
	}

	// Parse the JSON response
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Failed to unmarshal response: %v", err)
		return "", err
	}

	aiResponse := data.Choices[0].Message.Content

	return aiResponse, nil
}
