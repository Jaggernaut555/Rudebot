package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Constants
const (
	CmdChar = "~"
)

// CmdFuncType Command function type
type CmdFuncType func(*discordgo.Session, *discordgo.MessageCreate, []string)

// CmdFuncHelpType The type stored in the CmdFuncs map to map a function and helper text to a command
type CmdFuncHelpType struct {
	function          CmdFuncType
	help              string
	allowedChannelOnly bool
}

// CmdFuncsType The type of the CmdFuncs map
type CmdFuncsType map[string]CmdFuncHelpType

// CmdFuncs Commands to functions map
var CmdFuncs CmdFuncsType

// InitCmds Initializes the cmds map
func InitCmds() {
	CmdFuncs = CmdFuncsType{
		"help":     CmdFuncHelpType{cmdHelp, "Prints this list", false},
		"here":     CmdFuncHelpType{cmdHere, "Allows the bot to insult users in this channel", false},
		"nothere":	CmdFuncHelpType{cmdNotHere, "Restricts the bot from insulting users in this channel", true},
		"version":  CmdFuncHelpType{cmdVersion, "Outputs the current bot version", true},
		"insult":	CmdFuncHelpType{cmdInsult, "Follow by @user to insult that user", true},
		"rate":  CmdFuncHelpType{cmdRate, "Rate the insult dished out", true},
		"stats":    CmdFuncHelpType{cmdStats, "Displays stats about this bot", true},
	}
}

func HandleCommand(session *discordgo.Session, message *discordgo.MessageCreate, cmd string) {
	args := strings.Split(cmd, " ")
	if len(args) == 0 {
		return
	}
	CmdFuncHelpPair, ok := CmdFuncs[args[0]]

	if ok {
		if !CmdFuncHelpPair.allowedChannelOnly || isValidChannel(session, message.ChannelID) {
			CmdFuncHelpPair.function(session, message, args)
		}
	} else if isValidChannel(session, message.ChannelID) {
		var reply = fmt.Sprintf("I do not have command `%s`", args[0])
		session.ChannelMessageSend(message.ChannelID, reply)
	}
}

func cmdHelp(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	// Build array of the keys in CmdFuncs
	var keys []string
	for k := range CmdFuncs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build message (sorted by keys) of the commands
	var cmds = "Command notation: \n`" + CmdChar + "[command]`\n"
	cmds += "Commands:\n```\n"
	for _, key := range keys {
		cmds += fmt.Sprintf("%s - %s\n", key, CmdFuncs[key].help)
	}
	cmds += "```\n"
	session.ChannelMessageSend(message.ChannelID, cmds)
}

func cmdVersion(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	session.ChannelMessageSend(message.ChannelID, "Version: "+Version)
}

func cmdHere(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	Channels[message.ChannelID] = true
	var newArgs []string
	newArgs = append(newArgs, "insult")
	newArgs = append(newArgs, message.Author.Mention())

	cmdInsult(session, message, newArgs)
}

func cmdNotHere(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	Channels[message.ChannelID] = false
}

func cmdStats(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	var stats = "Stats:\n```\n"
	stats += fmt.Sprintf("Nouns: %d\n", NumNouns)
	stats += fmt.Sprintf("Adjectives: %d\n", NumAdjectives)
	stats += fmt.Sprintf("Adverbs: %d\n", NumAdverbs)
	stats += fmt.Sprintf("verbs: %d\n", NumVerbs)
	stats += "```"
	session.ChannelMessageSend(message.ChannelID, stats)
}

func cmdInsult(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		session.Channel(message.ChannelID)
		return
	}

	reply := NewInsult(args[1])
	session.ChannelMessageSend(message.ChannelID, reply)
}

func cmdRate(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {

}