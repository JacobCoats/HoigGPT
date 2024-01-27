package main

import (
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
)

// The way I'm storing conversation state on a per-user basis doesn't make sense. Thinking of storing conversation state up until
// a certain number of messages, wiping it after specific time period, or some combination. Also currently functionality is broken

type ConversationHandler struct {
	Conversation *genai.ChatSession
	StartTime    time.Time
	Model        *genai.GenerativeModel
}

var ConversationDuration = time.Minute * 60

func NewConversationHandler(model *genai.GenerativeModel) *ConversationHandler {
	return &ConversationHandler{
		Model: model,
	}
}

func (c *ConversationHandler) Get() *genai.ChatSession {
	// If there isn't an existing conversation or it's older than the duration, start a new one
	if c.Conversation == nil || time.Since(c.StartTime) > ConversationDuration {
		c.StartNewSession()
	}

	return c.Conversation
}

func (c *ConversationHandler) StartNewSession() {
	c.Conversation = c.Model.StartChat()
	c.StartTime = time.Now()
}

func (c *ConversationHandler) Print() {
	conversation := c.Conversation
	for _, message := range conversation.History {
		for _, part := range message.Parts {
			fmt.Println(part)
		}
	}
}
