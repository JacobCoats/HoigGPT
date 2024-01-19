package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
)

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()

	conversations = ConversationHandler{Conversations: make(map[string]*genai.ChatSession)}
	client, err := InitClient()
	if err != nil {
		fmt.Println("Error initializing request handler: ", err)
	}
	model = client.GenerativeModel("gemini-pro")
}

var token string
var conversations ConversationHandler
var model *genai.GenerativeModel

func main() {
	if token == "" {
		fmt.Println("No token provided")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating session: ", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening session: ", err)
	}

	// Create channel to listen for term signal. Quit when ctrl-c is pressed
	fmt.Println("Bot started successfully. Press ctrl-c to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!delete") {
		conversations.Delete(m.Author.ID)
		return
	}

	if strings.Contains(m.Content, "<@"+s.State.User.ID+">") {
		conversations.Add(m.Author.ID, model)
		conversations.Print(m.Author.ID)
	} else if m.ReferencedMessage == nil || m.ReferencedMessage.Author.ID != s.State.User.ID {
		return
	}

	var formattedMessage string
	// Remove the ping
	if strings.HasPrefix(m.Content, "<@"+s.State.User.ID+">") {
		formattedMessageWords := strings.Split(m.Content, " ")[1:]
		formattedMessage = strings.Join(formattedMessageWords, " ")
	} else {
		formattedMessage = m.Content
	}
	response, err := GetResponse(conversations.Get(m.Author.ID), formattedMessage)
	if err != nil {
		fmt.Println("Error retrieving resonse: ", err)
	}

	messageReference := discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
	}

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
