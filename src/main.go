package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/AB0529/dgo-wrapper/logger"
	dgowrapper "github.com/AB0529/dgo-wrapper/src"
	"github.com/AB0529/lady_luck/src/commands"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

type Conf struct {
	Token                           string
	Prefix                          string
	MessagesAppearInSpecificChannel bool
	MessageChannelID                string
}

var Config *Conf

func main() {
	f, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(f, &Config)

	s, err := dgowrapper.Initialize(&dgowrapper.Options{
		Prefixes: []string{Config.Prefix},
		Token:    Config.Token,
		Intent:   discordgo.IntentsAllWithoutPrivileged,
		Handlers: []interface{}{CreateMessage, Ready},
	})
	if err != nil {
		panic(err)
	}

	// db := DBInit()
	commands.Ping()
	commands.Flip()
	commands.Order()
	commands.Roll()

	dgowrapper.LogLoadedCommands(dgowrapper.Commands)
	// defer db.Close()

	err = dgowrapper.WaitForTerm(s)
	if err != nil {
		panic(err)
	}
}

func CreateMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Content) <= 1 {
		return
	}
	// Make sure message starts with prefix
	if !strings.HasPrefix(strings.ToLower(m.Message.Content), strings.ToLower(Config.Prefix)) {
		return
	}

	// Get the command
	cmd, err := dgowrapper.FindCommandOrAlias(Config.Prefix, m.Content)
	// Can't find command
	if err == dgowrapper.ErrNoCommandOrAliasFound {
		logger.Logf("WARN", "no command or alias found for %s", m.Content)
		return
	}

	// Make sure it's in guild
	channel, _ := s.Channel(m.ChannelID)
	if channel.Type == discordgo.ChannelTypeDM {
		logger.Logf("WARN", "channel type is DM")
		return
	}

    // Delete author message
	s.ChannelMessageDelete(m.ChannelID, m.Message.ID)

	var ctx *dgowrapper.Context

	if Config.MessagesAppearInSpecificChannel {
		// This is terrible. :)
		m.ChannelID = Config.MessageChannelID
		ctx = &dgowrapper.Context{
			Session: s,
			Message: m,
			Command: cmd,
			Prefix:  dgowrapper.Prefix,
		}
	} else {
		ctx = &dgowrapper.Context{
			Session: s,
			Message: m,
			Command: cmd,
			Prefix:  dgowrapper.Prefix,
		}
	}

	cmd.Handler(ctx)
}

func Ready(s *discordgo.Session, e *discordgo.Ready) {
	// Add mention prefix
	dgowrapper.Bot.Prefixes = append(dgowrapper.Bot.Prefixes, []string{
		fmt.Sprintf("<@%s> ", e.User.ID),
		fmt.Sprintf("<@!%s> ", e.User.ID),
	}...)

	logger.Logf("BOT", "%s#%s is ready", logger.Yellow.Sprint(e.User.Username), logger.Yellow.Sprint(e.User.Discriminator))
	logger.Logf("BOT", "Prefix is: %s", Config.Prefix)
	err := s.UpdateGameStatus(0, "Russian Roulette")
	logger.Die(err)
}

// func DBInit() *sql.DB {
//     db, err := sql.Open("postgres", Config.ConnStr)
//     if err != nil {
//         panic(err)
//     }
//
//     if err := db.Ping(); err != nil {
//         panic(err)
//     }
//     logger.Logf("DATABASE", "The database is connected")
//
//     return db
// }
