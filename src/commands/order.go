package commands

import (
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	dgowrapper "github.com/AB0529/dgo-wrapper/src"
	"github.com/AB0529/lady_luck/src/dice"
)

type mapEntry struct {
	key   string
	value int
}
func sortMapDescending(m map[string]int) map[string]int {
	// create a slice of map entries
	var entries []mapEntry
	for k, v := range m {
		entries = append(entries, mapEntry{k, v})
	}

	// sort the slice in descending order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].value > entries[j].value
	})

	// create a new map and add the sorted entries
	sortedMap := make(map[string]int)
	for _, entry := range entries {
		sortedMap[entry.key] = entry.value
	}

	return sortedMap
}

func Order() {
	dgowrapper.NewCommands([]*dgowrapper.Command{
		{
			Name:         "order",
			Aliases:      []string{"o"},
			Examples:     []string{"-order", "-o"},
			Descriptions: []string{"A way to create an order of battle from rolls."},
			Handler: func(ctx *dgowrapper.Context) {
				diceRollChanneID := "1047675722040361020"
				msg := strings.Join(strings.Split(strings.ToLower(ctx.Message.Content), " ")[1:], " ")

				if len(msg) < 1 {
					ctx.Embedf(":x: **Error** | Your command should look like this: `%sorder Mongi Idk` etc.", ctx.Prefix)
					return
				}

				names := strings.Split(msg, " ")
				result := map[string]int{}
				resultStr := ""

				rand.Seed(time.Now().UnixNano())

				for _, name := range names {
					// Filter out name from math exp
					re := regexp.MustCompile(`([a-zA-Z]+)([\+\-\*/])(\d+)`)
					matches := re.FindAllStringSubmatch(name, -1)

					if len(matches) <= 0 {
						_, diceNumber, mathExp, err := ParseDiceExpression("1d20")
						if err != nil {
							ctx.Embedf(":x: | Uh oh something went wrong...\n```css\n%s\n```", err.Error())
						}

						die := &dice.Dice{
							Sides:       diceNumber,
							Expressions: mathExp,
						}

						sum, _ := dice.RollDie(die)
						result[name] = sum
						continue
					}

					for _, match := range matches {
						symbol := match[2]
						number, _ := strconv.Atoi(match[3])

						_, diceNumber, mathExp, err := ParseDiceExpression(fmt.Sprintf("1d20%s%d", symbol, number))
						if err != nil {
							ctx.Embedf(":x: | Uh oh something went wrong...\n```css\n%s\n```", err.Error())
						}

						die := &dice.Dice{
							Sides:       diceNumber,
							Expressions: mathExp,
						}

						sum, _ := dice.RollDie(die)
						result[name] = sum
					}

				}

				result = sortMapDescending(result)
				i := 1

				for k, v := range result {
					resultStr += fmt.Sprintf("%d) %s -> %d\n", i, strings.ToUpper(k), v)
					i++
				}

				EmbedChannel(ctx, diceRollChanneID, fmt.Sprintf(":scroll: | Here is the final order!\n```css\n%s\n```", resultStr))
			},
		},
	})
}
