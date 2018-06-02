package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

var rollLead = "!roll"
var rollRegex = regexp.MustCompile(`!roll\s+(.*)$`)

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, rollLead) {
		r := doRoll(m.Content)
		s.ChannelMessageSend(m.ChannelID, r)
	}
}

func doRoll(in string) string {
	match := rollRegex.FindStringSubmatch(in)

	rollStr := ""
	if len(match) > 1 {
		rollStr = match[1]
	}

	// TODO this is dirty, need to clean it up A LOT
	rolls := ParseMany(rollStr)
	var buffer bytes.Buffer
	for i, dr := range rolls {
		if i > 1 {
			continue // TODO this is broken because of our parsing
		}
		if i > 0 {
			buffer.WriteString("\n")
		}

		roll := RollMany(dr.count, dr.dice)
		sum := 0

		buffer.WriteString(fmt.Sprintf("Rolling %dd%d ", dr.count, dr.dice))
		switch dr.oper {
		case "+":
			sum += dr.mod
			buffer.WriteString(fmt.Sprintf("+ %d ", dr.mod))
		case "-":
			sum -= dr.mod
			buffer.WriteString(fmt.Sprintf("- %d ", dr.mod))
		}

		buffer.WriteString("\t:\t[")
		for i, v := range roll {
			if i != 0 {
				buffer.WriteString(" + ")
			}
			buffer.WriteString(fmt.Sprintf("%d", v))
			sum += v
		}
		buffer.WriteString("] ")
		switch dr.oper {
		case "+":
			buffer.WriteString(fmt.Sprintf("+ %d ", dr.mod))
		case "-":
			buffer.WriteString(fmt.Sprintf("- %d ", dr.mod))
		}
		buffer.WriteString(fmt.Sprintf("= %d", sum))
	}
	return buffer.String()
}
