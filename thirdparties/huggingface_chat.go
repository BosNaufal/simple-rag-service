package thirdparties

import (
	app_constants "bos_personal_ai/env"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type HuggingFaceAIChatImpl struct {
}

func NewHuggingFaceAIChat() AIChatInterface {
	return &HuggingFaceAIChatImpl{}
}

func (e *HuggingFaceAIChatImpl) Prompt(systemPrompt string, userPrompt string, temp float32, maxToken int) (string, error) {
	var data CompletionResponse

	// url := "https://api.openai.com/v1/chat/completions"
	url := "https://api-inference.huggingface.co/models/meta-llama/Meta-Llama-3.1-70B-Instruct/v1/chat/completions"
	postData := map[string]interface{}{
		"model": "meta-llama/Meta-Llama-3.1-70B-Instruct",
		"messages": []map[string]interface{}{
			{
				"role":    "system",
				"content": systemPrompt,
			},
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
		"stream":      false,
		"temperature": temp,
		"max_tokens":  maxToken,
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
	// req.Header.Add("Authorization", "Bearer "+app_constants.OPENAI_API_KEY)
	req.Header.Add("Authorization", "Bearer "+app_constants.HUGGING_FACE_TOKEN)

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
