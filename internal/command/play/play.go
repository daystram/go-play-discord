package play

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tenntenn/goplayground"

	"github.com/daystram/go-play-discord/internal/server"
	"github.com/daystram/go-play-discord/internal/util"
)

const (
	runEmojiName = "🏗️"
	fmtEmojiName = "🧹"
)

func Register(srv *server.Server) error {
	srv.Session.AddHandler(listenReaction(srv))
	return nil
}

func listenReaction(srv *server.Server) func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	return func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		switch e.Emoji.Name {
		case runEmojiName:
			msg, err := s.ChannelMessage(e.ChannelID, e.MessageID)
			if err != nil {
				log.Println("command: play: run:", err)
				return
			}

			resultMsg, err := s.ChannelMessageSendReply(e.ChannelID, "_Building..._", msg.Reference())
			if err != nil {
				log.Println("command: play: run:", err)
				return
			}

			res, err := run(unwrapCode(msg.Content))
			if err != nil {
				log.Println("command: play: run:", err)
				return
			}

			go simulate(s, resultMsg, res)
		case fmtEmojiName:
			msg, err := s.ChannelMessage(e.ChannelID, e.MessageID)
			if err != nil {
				log.Println("command: play: fmt:", err)
				return
			}

			resultMsg, err := s.ChannelMessageSendReply(e.ChannelID, "_Formatting..._", msg.Reference())
			if err != nil {
				log.Println("command: play: fmt:", err)
				return
			}

			res, err := format(unwrapCode(msg.Content))
			if err != nil {
				log.Println("command: play: fmt:", err)
				return
			}

			if res.Error != "" {
				content := wrapResult("", "", util.TrimMessage(res.Error))
				_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID:      resultMsg.ID,
					Channel: resultMsg.ChannelID,
					Content: &content,
				})
				if err != nil {
					log.Println("command: play: fmt:", err)
				}

				finishResult(s, resultMsg, false)
			} else {
				content := wrapResult(res.Body, "", "")
				_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
					ID:      resultMsg.ID,
					Channel: resultMsg.ChannelID,
					Content: &content,
				})
				if err != nil {
					log.Println("command: play: fmt:", err)
				}

				finishResult(s, resultMsg, true)
			}
		}
	}
}

func run(code string) (*goplayground.RunResult, error) {
	c := &goplayground.Client{}
	res, err := c.Run(code)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func format(code string) (*goplayground.FormatResult, error) {
	c := &goplayground.Client{}
	res, err := c.Format(code, true)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func simulate(s *discordgo.Session, resultMsg *discordgo.Message, res *goplayground.RunResult) {
	var err error

	stdout := ""
	stderr := ""

	if res.Errors != "" {
		content := wrapResult("", "", util.TrimMessage(res.Errors))
		_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:      resultMsg.ID,
			Channel: resultMsg.ChannelID,
			Content: &content,
		})
		if err != nil {
			log.Println("command: play: simulate: failed simulating result: ", err)
			return
		}

		finishResult(s, resultMsg, false)
	} else {
		for _, e := range res.Events {
			time.Sleep(e.Delay)
			switch e.Kind {
			case "stdout":
				if strings.HasPrefix(e.Message, "\f") {
					stdout = ""
				}
				stdout = fmt.Sprint(stdout, e.Message)
			case "stderr":
				if strings.HasPrefix(e.Message, "\f") {
					stderr = ""
				}
				stderr = fmt.Sprint(stderr, e.Message)
			}

			content := wrapResult("", util.TrimMessage(stdout), util.TrimMessage(stderr))
			_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      resultMsg.ID,
				Channel: resultMsg.ChannelID,
				Content: &content,
			})
			if err != nil {
				log.Println("command: play: simulate: failed simulating result: ", err)
				return
			}
		}

		finishResult(s, resultMsg, true)
	}
}

func finishResult(s *discordgo.Session, msg *discordgo.Message, ok bool) {
	emoji := ""
	if ok {
		emoji = "✅"
	} else {
		emoji = "❌"
	}

	err := s.MessageReactionAdd(msg.ChannelID, msg.ID, emoji)
	if err != nil {
		log.Println("command: play: failed finishing result: ", err)
		return
	}
}

func unwrapCode(raw string) string {
	clean := raw
	clean = strings.TrimPrefix(clean, "```go")
	clean = strings.TrimSuffix(clean, "```")
	return clean
}

func wrapResult(code, stdout, stderr string) string {
	codeFmt := "```go\n%s\n```\n"
	outFmt := "`stdout`\n```\n%s\n```\n"
	errFmt := "`stderr`\n```\n%s\n```\n"

	wrapped := ""
	if code != "" {
		wrapped += fmt.Sprintf(codeFmt, code)
	}
	if stdout != "" {
		wrapped += fmt.Sprintf(outFmt, stdout)
	}
	if stderr != "" {
		wrapped += fmt.Sprintf(errFmt, stderr)
	}

	wrapped = strings.TrimSpace(wrapped)
	if wrapped == "" {
		wrapped = "_No output._"
	}

	return wrapped
}
