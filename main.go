package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var vcsession *discordgo.VoiceConnection
var buffer = make([][]byte, 0)

const (
	//TOKEN はBotのID
	TOKEN = "NzkyNzAyODkwMDQ0MjkzMTIw.X-hkGA.lFD_2baEvCltW7ceBjqurh6jLm8"
	//BotName はBotの名前
	BotName = "にいな"
)

func main() {

	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		fmt.Println("error:start\n", err)
		return
	}

	//on ready
	dg.AddHandler(ready)

	//on message
	dg.AddHandler(messageCreate)

	dg.AddHandler(loadSound)

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
	_ = s.UpdateStatus(0, "nil!!")
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
	fmt.Println("<< " + m.Content + " by " + nick)

	if m.Content == "こんにちわ" {
		s.ChannelMessageSend(m.ChannelID, "こんにちわ")
		fmt.Println("> こんにちわ")
	}
	if strings.Contains(m.Content, "ww") {
		s.ChannelMessageSend(m.ChannelID, "lol")
		fmt.Println("> lol")
	}

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
				err := playVoice1(s, g.ID, a.ChannelID)

				s.ChannelMessageSend(m.ChannelID, "にいながきたよ")
				if err != nil {
					fmt.Println("Error playing Voice: ", err)
				}
				return
			}
		}

	case strings.HasPrefix(m.Content, "!bye"):
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
				err := playVoice2(s, g.ID, a.ChannelID)

		        s.ChannelMessageSend(m.ChannelID, "帰るね。バイバーイ")
		        if err != nil {
			        fmt.Pritnln("")
		        }
				return
			}
		}

	if m.Content == m.Content {
		file, err := os.OpenFile("vc.txt", os.O_RDWR|os.O_CREATE, 0664)
		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()

		output := m.Content
		fmt.Fprintln(file, output)
		log.Printf("%T %v", output, output)

		cwd1, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}

		log.Println(cwd1)

		c := "open_jtalk"
		p := []string{
			"-x", "/usr/local/Cellar/open-jtalk/1.11/dic/",
			"-m", "/usr/local/Cellar/open-jtalk/1.11/voice/mei/mei_normal.htsvoice",
			"vc.txt",
			"-ow", "output.wav",
		}
		cmd1 := exec.Command(c, p...)
		log.Printf("%T %v", cmd1, cmd1)

		if err := cmd1.Run(); err != nil {
			log.Fatal("err1 Error: ", err)
		}
	}
}

func loadSound() error {

	file, err := os.Open("output.wav")
	if err != nil {
		fmt.Println("Error opening wav file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from wav file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from wav file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from wav file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

func playVoice1(s *discordgo.Session, guildID, channelID string) (err error) {

	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		return err
	}

	vc.Speaking(true)

	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	vc.Speaking(false)

	return nil
}

func playVoice2(s *discordgo.Session, guildID, channelID string) (err error) {

}