package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/TheDummyUser/goRss/config"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="
)

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

func SummarizeAi(c *fiber.Ctx, db *gorm.DB) error {
	var req struct {
		Context string `json:"context"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Context == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Context cannot be empty"})
	}

	prompt := fmt.Sprintf("Please summarize the following text in under 200 words, while maintaining key details. Provide sources if available. Format the output in markdown.:\n\n%s", req.Context)

	requestData := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to marshal request body: %v", err)})
	}

	reqHTTP, err := http.NewRequest("POST", geminiAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to create Gemini API request: %v", err)})
	}

	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("x-goog-api-key", config.Config("GEMINI_TOKEN")) // Or use an Authorization header

	client := &http.Client{}
	resp, err := client.Do(reqHTTP)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to send request to Gemini API: %v", err)})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return c.Status(resp.StatusCode).JSON(fiber.Map{"error": fmt.Sprintf("Gemini API request failed with status: %d, body: %s", resp.StatusCode, string(bodyBytes))})
	}

	var geminiResponse GeminiResponse
	err = json.NewDecoder(resp.Body).Decode(&geminiResponse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to decode Gemini API response: %v", err)})
	}

	var summary string
	if len(geminiResponse.Candidates) > 0 && len(geminiResponse.Candidates[0].Content.Parts) > 0 {
		summary = geminiResponse.Candidates[0].Content.Parts[0].Text
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No summary found in Gemini API response"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status_code": fiber.StatusOK,
		"summary":     summary,
	})
}
