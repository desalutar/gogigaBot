package gptApi

import (
	"R2D2/apps/gptAnswer/internal/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OpenAIClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient() *OpenAIClient {
	return &OpenAIClient{
		apiKey:     os.Getenv("GPT_KEY"),
		httpClient: &http.Client{},
	}
}

func (c *OpenAIClient) GetCompletion(ctx context.Context, prompt string) (string, error) {
	reqBody := model.ChatCompletionRequest{
		Model: os.Getenv("GPT_MODEL"),
		Messages: []model.ChatMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens: 100,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", os.Getenv("GPT_URL"), bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr model.OpenAIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Error.Message != "" {
			return "", fmt.Errorf("OpenAI error: %s", apiErr.Error.Message)
		}
		return "", fmt.Errorf("OpenAI request failed: status %d", resp.StatusCode)
	}

	var res model.ChatCompletionResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	if len(res.Choices) == 0 {
		return "", fmt.Errorf("empty response")
	}
	return res.Choices[0].Message.Content, nil
}
