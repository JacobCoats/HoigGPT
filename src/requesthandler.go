package main

import (
	"context"
	"fmt"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func InitClient() (*genai.Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyA_JSJbrNEm0_9pcVQEjNWcpsXzJ8kYclU"))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetResponse(chatSession *genai.ChatSession, message string) ([]genai.Part, error) {
	ctx := context.Background()

	fmt.Println("Requesting chat response for message: ", message)
	fmt.Println("History length: ", len(chatSession.History))
	res, err := chatSession.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return nil, err
	}

	return res.Candidates[0].Content.Parts, nil
}
