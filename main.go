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

func main2() {
	//a := "!roll"
	//b := "!roll "
	//c := "!roll x"
	//d := "!roll d20"

	//	fmt.Println(strings.HasPrefix(a))
	//	fmt.Println(strings.HasPrefix(b))
	//	fmt.Println(strings.HasPrefix(c))
	//	fmt.Println(strings.HasPrefix(d))

}

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

	fmt.Println(m.Content)
	if strings.HasPrefix(m.Content, rollLead) {
		fmt.Println("PREFIX " + m.Content)
		match := rollRegex.FindStringSubmatch(m.Content)
		rollStr := ""
		if len(match) > 1 {
			rollStr = match[1]
		}
		count, max := Parse(rollStr)
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("Rolling %dd%d (", count, max))
		o := RollMany(count, max)
		sum := 0
		for i, v := range o {
			if i != 0 {
				buffer.WriteString(" + ")
			}
			buffer.WriteString(fmt.Sprintf("%d", v))
			sum += v
		}
		buffer.WriteString(fmt.Sprintf(") = %d", sum))
		s.ChannelMessageSend(m.ChannelID, buffer.String())
	}
}
