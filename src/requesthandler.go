package main

import (
	"context"
	"fmt"

	genai "github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func InitClient(token string) (*genai.Client, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(token))
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
