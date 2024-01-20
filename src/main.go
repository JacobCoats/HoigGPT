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
	flag.StringVar(&discordToken, "d", "", "Discord Token")
	flag.StringVar(&bardToken, "b", "", "Bard Token")
	flag.Parse()

	if bardToken == "" {
		fmt.Println("No Bard token provided")
		return
	}

	conversations = ConversationHandler{Conversations: make(map[string]*genai.ChatSession)}
	client, err := InitClient(bardToken)
	if err != nil {
		fmt.Println("Error initializing request handler: ", err)
	}
	model = client.GenerativeModel("gemini-pro")
}

var discordToken string
var bardToken string
var conversations ConversationHandler
var model *genai.GenerativeModel

func main() {
	if bardToken == "" {
		return
	}

	if discordToken == "" {
		fmt.Println("No Discord token provided")
		return
	}

	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("Error creating session: ", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening session: ", err)
		return
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

	HandleMessage(s, m)
}
