package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type AISuggestionRequest struct {
	Prompt string `json:"prompt"`
}

type AISuggestionResponse struct {
	Suggestions []string `json:"suggestions"`
}

func GetAITaskSuggestions(prompt string) ([]string, error) {
	// Example: Call OpenAI/Gemini API (replace with actual API endpoint and key)
	url := "https://api.openai.com/v1/completions" // Replace with actual API endpoint
	apiKey := "your-openai-api-key"               // Replace with your OpenAI/Gemini API key from config

	requestBody := AISuggestionRequest{Prompt: prompt}
	jsonBody, _ := json.Marshal(requestBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response AISuggestionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Suggestions, nil
}