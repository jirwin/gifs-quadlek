package gifs

import (
	"context"

	"fmt"

	"strings"

	"net/url"

	"github.com/jirwin/quadlek/quadlek"
)

var gifs *Gifs

func gifCommand(ctx context.Context, cmdChannel <-chan *quadlek.CommandMsg) {
	for {
		select {
		case cmdMsg := <-cmdChannel:
			text := cmdMsg.Command.Text
			if text != "" {
				cmdMsg.Store.Get(text, func(v []byte) error {
					if v != nil {
						cmdMsg.Command.Reply() <- &quadlek.CommandResp{
							Text:      string(v),
							InChannel: true,
						}
						return nil
					}
					r, err := gifs.Translate(text)
					if err != nil {
						cmdMsg.Command.Reply() <- &quadlek.CommandResp{
							Text:      fmt.Sprintf("an error occured: %s", err.Error()),
							InChannel: false,
						}
						return nil
					}
					cmdMsg.Command.Reply() <- &quadlek.CommandResp{
						Text:      r,
						InChannel: true,
					}
					return nil
				})
			}

		case <-ctx.Done():
			return
		}
	}
}

func gifSaveCommand(ctx context.Context, cmdChannel <-chan *quadlek.CommandMsg) {
	for {
		select {
		case cmdMsg := <-cmdChannel:
			text := cmdMsg.Command.Text
			parts := strings.Split(text, " ")
			if len(parts) < 2 {
				cmdMsg.Command.Reply() <- &quadlek.CommandResp{
					Text:      "Malformed command: /gsave <url> phrase to save",
					InChannel: false,
				}
				continue
			}
			gUrl, err := url.Parse(parts[1])
			if err != nil {
				cmdMsg.Command.Reply() <- &quadlek.CommandResp{
					Text:      fmt.Sprintf("Invalid url: %s", parts[1]),
					InChannel: false,
				}
				continue
			}

			phrase := strings.Join(parts[2:], " ")

			err = cmdMsg.Store.Update(phrase, []byte(gUrl.String()))
			if err != nil {
				cmdMsg.Command.Reply() <- &quadlek.CommandResp{
					Text:      fmt.Sprintf("Unable to save phrase: %s", err.Error()),
					InChannel: false,
				}
				continue
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
			quadlek.MakeCommand("gsave", gifSaveCommand),
		},
		nil,
		nil,
		nil,
		nil,
	)
}
