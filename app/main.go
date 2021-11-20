package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"

	"discord-bot/src/config"
)

var vcsession *discordgo.VoiceConnection
var buffer = make([][]byte, 0)

func main() {

	config := config.NewConfig()

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error: start\n", err)
		return
	}

	//on ready
	dg.AddHandler(ready)

	//on message
	dg.AddHandler(messageCreate)

	//ギルド（チャンネルを含む）に関する情報が必要です、メッセージと音声状態。
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates)

	//websocketを開いてRunning開始
	err = dg.Open()
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Println("BOT Running...")

	//シグナル受け取り可にしてチャネル受け取りを待つ（受け取ったら終了）
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc //待機する
	dg.Close()
}

//Botの状態表示
func ready(s *discordgo.Session, event *discordgo.Ready) {

	s.UpdateGameStatus(0, "Listening!!")

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	nick := m.Author.Username
	member, err := s.State.Member(m.GuildID, m.Author.ID)
	if err == nil && member.Nick != "" {
		nick = member.Nick
	}
	fmt.Println(m.Author.Username)

	fmt.Println("<< " + m.Content + " by " + nick)

	switch {
	case strings.HasPrefix(m.Content, "!join"):
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			log.Println(err)
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			log.Println(err)
		}

		for _, a := range g.VoiceStates {
			if a.UserID == m.Author.ID {
				vcsession, _ = s.ChannelVoiceJoin(c.GuildID, a.ChannelID, false, false)
				if err != nil {
					fmt.Println(err)
				}

				s.ChannelMessageSend(m.ChannelID, "にいながきたよ")

			}
		}

	case strings.HasPrefix(m.Content, "!bye"):

		// vcsession.Speaking(false)

		s.ChannelMessageSend(m.ChannelID, "帰るね。バイバーイ")
		vcsession.Disconnect()
	}

	if m.Content == "こんにちわ" {
		CreateWav(vcsession, m)

		fileName := "output.wav"

		fmt.Println("reading file name: ", fileName)

		vcsession.Speaking(true)

		// vc, err := s.ChannelVoiceJoin(m.GuildID, m.ChannelID, false, false)
		// if err != nil {
		// 	fmt.Println("discord voice Connection:", err)
		// }
		// fmt.Println()
		dgvoice.PlayAudioFile(vcsession, fileName, make(chan bool))

		err := vcsession.Speaking(false)
		if err != nil {
			fmt.Println("Speaking false error: ", err)
		}

		s.ChannelMessageSend(m.ChannelID, "こんにちわ")
		fmt.Println("> こんにちわ")
	}
	if strings.Contains(m.Content, "ww") {
		s.ChannelMessageSend(m.ChannelID, "lol")
		fmt.Println("> lol")
	}
}

// //CreateWav ここでwav音声ファイルを作成する
func CreateWav(v *discordgo.VoiceConnection, m *discordgo.MessageCreate) {

	if m.Content == m.Content {
		file, err := os.OpenFile("vc.txt", os.O_RDWR|os.O_CREATE, 0664)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		output := m.Content
		fmt.Fprintln(file, output)
		log.Printf("%T %v", output, output)

		c := "open_jtalk"
		p := []string{
			"-x", "discord-bot/open-jtalk/1.11/dic/",
			"-m", "discord-bot/open-jtalk/1.11/voice/mei/mei_normal.htsvoice",
			"vc.txt",
			"-ow", "output.wav",
		}
		cmd := exec.Command(c, p...)

		cmd.Run()

		file, err = os.Open("output.wav")

		if err != nil {
			log.Println("Error opening file: ", err)
		}
	}

}
