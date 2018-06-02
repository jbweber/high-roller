// From discordgo/examples/pingpong/main.go commit/abc85a2de0d321c21196a3f4374db5d1bcc9c9ea
// Copyright (c) 2018, Jeff Weber
// Copyright (c) 2015, Bruce Marriner
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of discordgo nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
	totals := make([]int, len(rolls))
	for i, dr := range rolls {
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
		totals[i] = sum
	}

	if len(rolls) > 1 {
		total := 0
		buffer.WriteString("\n")
		buffer.WriteString("Total Sum ")
		for i, v := range totals {
			if i != 0 {
				buffer.WriteString(" + ")
			}

			buffer.WriteString(fmt.Sprintf("%d", v))
			total += v
		}

		buffer.WriteString(fmt.Sprintf(" = %d", total))
	}

	return buffer.String()
}
