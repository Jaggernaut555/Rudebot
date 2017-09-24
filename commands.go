package main

import (
	"fmt"
	"math/rand"
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
	function           CmdFuncType
	help               string
	allowedChannelOnly bool
}

// CmdFuncsType The type of the CmdFuncs map
type CmdFuncsType map[string]CmdFuncHelpType

// CmdFuncs Commands to functions map
var CmdFuncs CmdFuncsType

// InitCmds Initializes the cmds map
func InitCmds() {
	CmdFuncs = CmdFuncsType{
		"help":    CmdFuncHelpType{cmdHelp, "Prints this list", false},
		"here":    CmdFuncHelpType{cmdHere, "Allows the bot to insult users in this channel", false},
		"nothere": CmdFuncHelpType{cmdNotHere, "Restricts the bot from insulting users in this channel", true},
		"version": CmdFuncHelpType{cmdVersion, "Outputs the current bot version", true},
		"insult":  CmdFuncHelpType{cmdInsult, "Insults user mentioned in [arguments]. Leave [arguments] blank to insult a random user in the server", true},
		"rate":    CmdFuncHelpType{cmdRate, "Rate the insult dished out. Use 'up', 'down', 'trash', or 'lmao'", true},
		"stats":   CmdFuncHelpType{cmdStats, "Displays stats about this bot", true},
		"define":  CmdFuncHelpType{cmdDefine, "Displays definition of [arguments]", true},
		"best":    CmdFuncHelpType{cmdBest, "~insult that selects the highest rated", true},
		"worst":   CmdFuncHelpType{cmdWorst, "~insult that selects the lowest rated", true},
		"good":    CmdFuncHelpType{cmdGood, "~insult that selects only positive rate", true},
		"bad":     CmdFuncHelpType{cmdBad, "~insult that selects only negative rated", true},
		"last":    CmdFuncHelpType{cmdLast, "~insult that selects the last used insult", true},
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
		SendReply(session, message, reply)
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
	var cmds = "Command notation: \n`" + CmdChar + "[command] [arguments]`\n"
	cmds += "Commands:\n```\n"
	for _, key := range keys {
		cmds += fmt.Sprintf("%s - %s\n", key, CmdFuncs[key].help)
	}
	cmds += "```\n"
	SendReply(session, message, cmds)
}

func cmdVersion(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	SendReply(session, message, "Version: "+Version)
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
	SendReply(session, message, stats)
}

func cmdRate(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		SendReply(session, message, "No rating given")
		return
	}
	switch args[1] {
	case "lmao":
		Rate(2)
	case "up":
		Rate(1)
	case "down":
		Rate(-1)
	case "trash":
		Rate(-2)
	default:
		SendReply(session, message, "No Rating Given")
	}
}

func cmdDefine(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		SendReply(session, message, "No search term given")
		return
	}

	definitions := DefineWord(args[1])

	SendReply(session, message, definitions)
}

func validateInsult(session *discordgo.Session, message *discordgo.MessageCreate, args []string) []string {
	if len(args) < 2 {
		var channel, err = session.Channel(message.ChannelID)
		if err != nil {
			fmt.Printf("Could not find channel, %s\n", err)
			return nil
		}

		guild, err := session.Guild(channel.GuildID)

		members := guild.Members
		user := members[rand.Intn(len(members))].User

		args = append(args, user.Mention())
	}

	return args
}

func cmdInsult(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	args = validateInsult(session, message, args)
	if args == nil {
		fmt.Printf("Could not create valid insult")
		return
	}
	reply := RandomInsult(args[1])
	SendReply(session, message, reply)
}

func cmdBest(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	args = validateInsult(session, message, args)
	if args == nil {
		fmt.Printf("Could not create valid insult")
		return
	}
	reply := BestInsult(args[1])
	SendReply(session, message, reply)
}

func cmdWorst(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	args = validateInsult(session, message, args)
	if args == nil {
		fmt.Printf("Could not create valid insult")
		return
	}
	reply := WorstInsult(args[1])
	SendReply(session, message, reply)
}

func cmdGood(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	args = validateInsult(session, message, args)
	if args == nil {
		fmt.Printf("Could not create valid insult")
		return
	}
	reply := GoodInsult(args[1])
	SendReply(session, message, reply)
}

func cmdBad(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	args = validateInsult(session, message, args)
	if args == nil {
		fmt.Printf("Could not create valid insult")
		return
	}
	reply := BadInsult(args[1])
	SendReply(session, message, reply)
}

func cmdLast(session *discordgo.Session, message *discordgo.MessageCreate, args []string) {
	args = validateInsult(session, message, args)
	if args == nil {
		fmt.Printf("Could not create valid insult")
		return
	}
	reply := LastInsult(args[1])
	SendReply(session, message, reply)
}
