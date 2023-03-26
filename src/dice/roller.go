package dice

import (
	"fmt"
	"math/rand"
	"time"
)

// PerformMath perform math on the total
func PerformMath(total int, die *Dice) (int, string) {
	operationsPerfomed := ""

	for operation, num := range die.Expressions {
		switch operation {
		case "+":
			total += num
			operationsPerfomed += fmt.Sprintf("+ %d\n", num)
			break
		case "-":
			total -= num
			operationsPerfomed += fmt.Sprintf("- %d\n", num)
			break
		case "*":
			total *= num
			operationsPerfomed += fmt.Sprintf("* %d\n", num)
			break
		case "/":
			total /= num
			operationsPerfomed += fmt.Sprintf("/ %d\n", num)
			break
		}
	}

	return total, operationsPerfomed
}

// RollDie rolls a single die
func RollDie(die *Dice) (int, string) {
	total := RandNumber(1, die.Sides)
	operationsPerfomed := ""

	if len(die.Expressions) <= 0 {
		return RandNumber(1, die.Sides), ""
	}

	return total, operationsPerfomed
}

// RollDice rolls a slice of dice
func RollDice(dice []*Dice) ([]int, string) {
	result := make([]int, len(dice))
	ops := ""

	for _, die := range dice {
		total, opsPer := RollDie(die)
		result = append(result, total)
		ops += opsPer
	}

	return result, ops
}

// RandNumber generates a random number from min and max values
func RandNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
