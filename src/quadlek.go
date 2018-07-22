package gifs

import (
	"context"

	"fmt"

	"github.com/jirwin/quadlek/quadlek"
)

var gifs *Gifs

func gifCommand(ctx context.Context, cmdChannel <-chan *quadlek.CommandMsg) {
	for {
		select {
		case cmdMsg := <-cmdChannel:
			text := cmdMsg.Command.Text
			if text != "" {
				r, err := gifs.Translate(text)
				if err != nil {
					cmdMsg.Command.Reply() <- &quadlek.CommandResp{
						Text:      fmt.Sprintf("an error occured: %s", err.Error()),
						InChannel: false,
					}
					continue
				}
				cmdMsg.Command.Reply() <- &quadlek.CommandResp{
					Text:      r,
					InChannel: true,
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

func Register(apiKey string) quadlek.Plugin {
	gifs = NewGifs(apiKey, "PG-13")
	return quadlek.MakePlugin(
		"gifs",
		[]quadlek.Command{
			quadlek.MakeCommand("g", gifCommand),
		},
		nil,
		nil,
		nil,
		nil,
	)
}
