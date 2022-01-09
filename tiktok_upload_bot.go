package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	targetChannelID := os.Getenv("DISCORD_TARGET_CHANNEL_ID")
	session, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		fmt.Printf("unable to create discord bot session, see: %s", err)
		return
	}

	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.MessageReaction.Emoji.Name == "✅" && m.ChannelID == targetChannelID {
			msg, err := session.ChannelMessage(targetChannelID, m.MessageID)
			if err != nil {
				session.ChannelMessageSend(targetChannelID, fmt.Sprintf("something went wrong, with message retrieval, see: %s", err))
			} else {
				if len(msg.Attachments) != 1 {
					session.ChannelMessageSend(targetChannelID, "Please ensure you only attach on file per upload.")
					return
				}
				url := msg.Attachments[0].URL
				session.ChannelMessageSend(targetChannelID, fmt.Sprintf("url of file to upload: %s", url))
				resp, err := http.Get(url)
				if err != nil {
					session.ChannelMessageSend(targetChannelID, "Failed to retrieve blob from CDN. Exiting upload process.")
					return
				}
				respBlob, err := io.ReadAll(resp.Body)
				resp.Body.Close()
				session.ChannelMessageSend(targetChannelID, fmt.Sprintf("blob buffered, size: %d", len(respBlob)))
			}
		}
	})

	err = session.Open()
	if err != nil {
		fmt.Printf("unable to open discord bot session, see: %s", err)
		return
	}

	fmt.Println("Bot running... Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	session.Close()
}
