package command

import (
	"github.com/daystram/go-play-discord/internal/server"
)

type RegisterFunc func(srv *server.Server) error

func commands() []RegisterFunc {
	return []RegisterFunc{}
}

func RegisterAll(srv *server.Server) error {
	for _, f := range commands() {
		err := f(srv)
		if err != nil {
			return err
		}
	}

	return nil
}
