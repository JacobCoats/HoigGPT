package main

import (
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

// The way I'm storing conversation state on a per-user basis doesn't make sense. Thinking of storing conversation state up until
// a certain number of messages, wiping it after specific time period, or some combination. Also currently functionality is broken

type ConversationHandler struct {
	Conversations map[string]*genai.ChatSession
}

func (c ConversationHandler) Get(userID string) *genai.ChatSession {
	return c.Conversations[userID]
}

func (c ConversationHandler) Add(userID string, model *genai.GenerativeModel) {
	c.Conversations[userID] = model.StartChat()
}

func (c ConversationHandler) Delete(userID string) {
	delete(c.Conversations, userID)
}

func (c ConversationHandler) Print(userID string) {
	conversation := c.Conversations[userID]
	for _, message := range conversation.History {
		for _, part := range message.Parts {
			fmt.Println(part)
		}
	}
}
