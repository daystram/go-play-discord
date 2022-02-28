package server

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"github.com/daystram/go-play-discord/internal/config"
)

type Server struct {
	Session *discordgo.Session
	Config  *config.Config
}

func Start(cfg *config.Config) (*Server, error) {
	s, err := discordgo.New(fmt.Sprintf("Bot %s", cfg.BotToken))
	if err != nil {
		return nil, err
	}

	err = s.Open()
	if err != nil {
		return nil, err
	}

	return &Server{
		Session: s,
		Config:  cfg,
	}, nil
}

func (srv *Server) Stop() error {
	return srv.Session.Close()
}
