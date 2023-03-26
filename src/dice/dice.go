package dice

// Expression the math expressions to apply to a rolled dice
type Expression map[string]int

// Dice a representation of a dice, each dice has n number of sides
// and expressions to perform on the result
type Dice struct {
	Sides       int
	Expressions Expression
}
