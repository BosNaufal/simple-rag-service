package thirdparties

import (
	app_constants "bos_personal_ai/env"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type EmbeddingThirdPartyInterface interface {
	GetEmbeddingFromString(content string) (string, error)
}

type EmbeddingThirdPartyImpl struct {
}

type OpenAIEmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string          `json:"object"`
		Index     int             `json:"index"`
		Embedding json.RawMessage `json:"embedding"`
	} `json:"data"`
	EmbeddingString string
}

func NewEmbeddingOpenAIEmbedding() EmbeddingThirdPartyInterface {
	return &EmbeddingThirdPartyImpl{}
}

func (e *EmbeddingThirdPartyImpl) GetEmbeddingFromString(content string) (string, error) {
	var data OpenAIEmbeddingResponse

	// Define the URL of the third-party API
	url := "https://api.openai.com/v1/embeddings"

	postData := map[string]string{
		"input": content,
		"model": "text-embedding-3-small",
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
	json.Unmarshal(body, &data)
	embeddingString := string(data.Data[0].Embedding)

	return embeddingString, nil
}
