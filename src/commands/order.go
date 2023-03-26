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

func SortMapByValue(inputMap map[string]int) map[string]int {
	// Convert the map to a slice of key-value pairs
	pairs := make([][2]interface{}, len(inputMap))
	i := 0
	for k, v := range inputMap {
		pairs[i] = [2]interface{}{k, v}
		i++
	}

	// Sort the slice by value
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i][1].(int) > pairs[j][1].(int)
	})

	// Convert the sorted slice back to a map
	outputMap := make(map[string]int, len(pairs))
	for _, pair := range pairs {
		key := pair[0].(string)
		value := pair[1].(int)
		outputMap[key] = value
	}
	return outputMap
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

				result = SortMapByValue(result)
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
