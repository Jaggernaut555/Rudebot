package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Constants
const (
	Version = "v0.1.1"
)

// Global vars
var (
	Token    string
	Channels map[string]bool
)

func init() {
	fmt.Println("Rudebot lives...")
	flag.StringVar(&Token, "t", "", "Discord Authentication Token")
	flag.Parse()

	Channels = map[string]bool{}
	InitInsults()
	//InitRatings()
	InitCmds()
}

func main() {
	if Token == "" {
		fmt.Println("You must provide a Discord authentication token (-t)")
		return
	}

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// add a handler for when messages are posted
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Do not talk to self
	if message.Author.ID == session.State.User.ID {
		return
	}

	if strings.HasPrefix(message.Content, CmdChar) {
		HandleCommand(session, message, strings.TrimPrefix(message.Content, CmdChar))
	}
}

func isValidChannel(session *discordgo.Session, channelID string) bool {
	return Channels[channelID]
}