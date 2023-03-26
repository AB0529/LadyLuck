package commands

import dgowrapper "github.com/AB0529/dgo-wrapper/src"

func Ping() {
	dgowrapper.NewCommands([]*dgowrapper.Command{
		{
			Name:         "ping",
			Aliases:      []string{"pong"},
			Examples:     []string{"-ping"},
			Descriptions: []string{"standard ping-pong command"},
			Handler: func(ctx *dgowrapper.Context) {
				ctx.Send("Pong")
			},
		},
	})
}
