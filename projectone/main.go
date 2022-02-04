package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Post("/topwords", func(c *fiber.Ctx) error {
		var req TopWordRequest
		c.BodyParser(&req)
		topWords, sessionID := GetTopWords(10, req.SessionID, req.Text)
		result := TopWordResponse{
			TopWords:  topWords,
			SessionID: sessionID,
		}
		return c.JSON(result)
	})
	app.Listen(":8080")
}

type TopWordRequest struct {
	Text      string `json:"text"`
	SessionID string `json:"session_id"`
}

type TopWordResponse struct {
	TopWords  []TopWord `json:"top_words"`
	SessionID string    `json:"session_id"`
}
