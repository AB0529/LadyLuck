package commands

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	dgowrapper "github.com/AB0529/dgo-wrapper/src"
	"github.com/AB0529/lady_luck/src/dice"
	"github.com/bwmarrin/discordgo"
)

func ParseDiceExpression(expression string) (rollCount int, diceNumber int, mathExprs dice.Expression, err error) {
	// Initialize the map of math expressions
	mathExprs = make(map[string]int)

	// Use a regular expression to match the roll count, dice number, and math expressions
	re := regexp.MustCompile(`^\s*(\d+)\s*d\s*(\d+)\s*(.*)\s*$`)
	match := re.FindStringSubmatch(expression)

	// Check if the expression does not match the format
	if match == nil {
		err = fmt.Errorf("expression does not match format")
		return
	}

	// Extract the roll count and dice number from the match
	rollCount, err = strconv.Atoi(match[1])
	if err != nil {
		return
	}
	diceNumber, err = strconv.Atoi(match[2])
	if err != nil {
		return
	}

	// Extract the math expressions from the match
	mathExprStr := strings.TrimSpace(match[3])
	for len(mathExprStr) > 0 {
		// Use a regular expression to match the next math expression and number
		re = regexp.MustCompile(`^([\+\-\*/])\s*(\d+)`)
		match = re.FindStringSubmatch(mathExprStr)

		// Check if the math expression was found
		if match == nil {
			err = fmt.Errorf("invalid math expression")
			return
		}

		// Extract the math expression and number from the match
		mathExpr := match[1]
		mathNumber, err := strconv.Atoi(match[2])
		if err != nil {
			return 0, 0, nil, err
		}

		// Add the math expression and number to the map
		mathExprs[mathExpr] = mathNumber

		// Remove the matched math expression and number from the math expression string
		mathExprStr = strings.TrimSpace(strings.TrimPrefix(mathExprStr, match[0]))
	}

	return
}

func Roll() {
	dgowrapper.NewCommands([]*dgowrapper.Command{
		{
			Name:         "roll",
			Aliases:      []string{"r"},
			Examples:     []string{"!roll 1d20"},
			Descriptions: []string{"roll dice"},
			Handler: func(ctx *dgowrapper.Context) {
				msg := strings.Join(strings.Split(strings.ToLower(ctx.Message.Content), " ")[1:], " ")

				if len(msg) < 1 {
					ctx.Embed(":x: | No dice roll specified. `Ex. 1d20, 2d4`")
					return
				}

				rollCount, diceNumber, mathExp, err := ParseDiceExpression(msg)
				if err != nil {
					ctx.Embed(":x: Error | Your command should look like: `1d20, 2d20-3, 1d4+2, etc.` Try again.")
					return
				}

				// Create the dice struct
				die := &dice.Dice{
					Sides:       diceNumber,
					Expressions: mathExp,
				}

				rolls := []int{}
				stringRolls := make([]string, len(rolls))

				opsStr := ""
				if rollCount == 1 {

					t, _ := dice.RollDie(die)
					rolls = append(rolls, t)
				} else {
					for i := 0; i < rollCount; i++ {
						t, _ := dice.RollDie(die)
						rolls = append(rolls, t)
					}
				}

				for _, n := range rolls {
                    stringRolls = append(stringRolls, strconv.Itoa(n))
				}
				
				// Generate a roll based off file, owner only
				content, err := ioutil.ReadFile("./r.txt")
				if err != nil {
					panic(err)
				}

                // TODO: Use config file for ID
				if ctx.Message.Author.ID == "184157133187710977" && string(content) != "" {
					cnt := strings.TrimSpace(string(content))
					ctx.Embedf(fmt.Sprintf(":game_die: %s | <@%s> rolled a **%s**!", msg, ctx.Message.Author.ID, cnt))

					if cnt == "1" {
						ctx.Send("L")
					}
					return
				}

				// Generate SQL statement
				// values := []string{}
				//
				// for _, roll := range rolls {
				//     r, _ := strconv.Atoi(roll)
				//     values = append(values, fmt.Sprintf("('%s', 'd%d', %d)", ctx.Message.Author.ID, diceNumber, r))
				// }
				//
				// _, err = db.Exec(fmt.Sprintf(`INSERT INTO rolls("user", die, roll) VALUES %s`, strings.Join(values, ",")))

				if err != nil {
					ctx.Embedf(":x: | Something horrible has gone wrong, check logs.")
					fmt.Println(err)
					return
				}

				if rollCount <= 1 {
					ctx.Embedf(fmt.Sprintf(":game_die: %s | <@%s> rolled a **%d**!", msg, ctx.Message.Author.ID, rolls[0]))
					return
				}

				sum := 0
				for _, i := range rolls {
					sum += i
				}
                sum, opsStr = dice.PerformMath(sum, die)

				ctx.Embedf(fmt.Sprintf(":game_die: %s | <@%s>\n```css\n%s\n%s\n=%d```", msg, ctx.Message.Author.ID, strings.Join(stringRolls, "+"), opsStr, sum))
			},
		},
	})
}

func EmbedChannel(ctx *dgowrapper.Context, ID string, content string) {
	ctx.Session.ChannelMessageSendComplex(ID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color:       rand.Intn(10000000),
			Description: content,
		},
	})
}
