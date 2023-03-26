package commands

import (
	"fmt"
	"math/rand"
	"time"

	dgowrapper "github.com/AB0529/dgo-wrapper/src"
	"github.com/bwmarrin/discordgo"
)

func Flip() {
	dgowrapper.NewCommands([]*dgowrapper.Command{
		{
			Name:         "flip",
			Aliases:      []string{"f"},
			Examples:     []string{"-flip"},
			Descriptions: []string{"flips a coin"},
			Handler: func(ctx *dgowrapper.Context) {
				rand.Seed(time.Now().UnixNano())
				c := rand.Intn(10)
				
				if c%2 == 0 {
					ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
						Embed: &discordgo.MessageEmbed{
							Color:       rand.Intn(10000000),
							Description: fmt.Sprintf(":coin: | <@%s> flipped a coin!", ctx.Message.Author.ID),
							Image: &discordgo.MessageEmbedImage{
								URL:    "http://www.clker.com/cliparts/7/d/e/0/139362185558690588heads-hi.png",
								Width:  512,
								Height: 512,
							},
						},
					})
					return
				}
				ctx.Session.ChannelMessageSendComplex(ctx.Message.ChannelID, &discordgo.MessageSend{
					Embed: &discordgo.MessageEmbed{
						Color:       rand.Intn(10000000),
						Description: fmt.Sprintf(":coin: | <@%s> flipped a coin!", ctx.Message.Author.ID),
						Image: &discordgo.MessageEmbedImage{
							URL:    "https://media.istockphoto.com/id/476142091/photo/quarter-dollar-us-coin-isolated-on-white.jpg?s=612x612&w=0&k=20&c=wNzr7m0Z3dhlf8_O1G3EFNz8u2tALVobVs4K4XfFN5c=",
							Width:  512,
							Height: 512,
						},
					},
				})
			},
		},
	})
}
