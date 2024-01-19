package main

import (
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

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
