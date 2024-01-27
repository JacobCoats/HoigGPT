package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate, promptModifier string) {
	messageReference := discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
	}

	// Remove the ping for the query
	var formattedQuery string
	if strings.HasPrefix(m.Content, "<@"+s.State.User.ID+">") {
		formattedQueryWords := strings.Split(m.Content, " ")[1:]
		formattedQuery = strings.Join(formattedQueryWords, " ")
	} else {
		formattedQuery = m.Content
	}

	var builder strings.Builder
	builder.WriteString(promptModifier)
	builder.WriteString(formattedQuery)
	modifiedQuery := builder.String()

	response, err := GetResponse(conversations.Get(), modifiedQuery)
	if err != nil {
		fmt.Println("Error retrieving resonse: ", err)

		if strings.Contains(err.Error(), "blocked: prompt") || strings.Contains(err.Error(), "blocked: candidate") {
			errorMsg := &discordgo.MessageSend{
				Content:         "Sorry, I can't respond to that. Please ridicule me",
				Reference:       &messageReference,
				AllowedMentions: &discordgo.MessageAllowedMentions{RepliedUser: true},
			}

			s.ChannelMessageSendComplex(m.ChannelID, errorMsg)
		}

		return
	}

	// Convert the genai Part object returned by API to string
	var responseString string
	for _, part := range response {
		partString := fmt.Sprintf("%d", part)
		partString = partString[15 : len(partString)-1]
		responseString += partString
	}

	msgSend := &discordgo.MessageSend{
		Content:         responseString,
		Reference:       &messageReference,
		AllowedMentions: &discordgo.MessageAllowedMentions{RepliedUser: true},
	}

	s.ChannelMessageSendComplex(m.ChannelID, msgSend)
}
