package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	conversations.Add(m.Author.ID, model)

	// Remove the ping for the query
	var formattedQuery string
	if strings.HasPrefix(m.Content, "<@"+s.State.User.ID+">") {
		formattedQueryWords := strings.Split(m.Content, " ")[1:]
		formattedQuery = strings.Join(formattedQueryWords, " ")
	} else {
		formattedQuery = m.Content
	}
	response, err := GetResponse(conversations.Get(m.Author.ID), formattedQuery)
	if err != nil {
		fmt.Println("Error retrieving resonse: ", err)
	}

	// Convert the genai Part object returned by API to string
	var responseString string
	for _, part := range response {
		partString := fmt.Sprintf("%d", part)
		partString = partString[15 : len(partString)-1]
		responseString += partString
	}

	messageReference := discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
	}

	msgSend := &discordgo.MessageSend{
		Content:         responseString,
		Reference:       &messageReference,
		AllowedMentions: &discordgo.MessageAllowedMentions{RepliedUser: true},
	}

	s.ChannelMessageSendComplex(m.ChannelID, msgSend)
}
