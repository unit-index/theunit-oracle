package marshal

import "strings"

// colorCode represents ANSII escape code for color formatting.
type colorCode string

const (
	reset colorCode = "\033[0m"
	red   colorCode = "\033[31m"
	green colorCode = "\033[32m"
)

var colorsEnabled = true

// disableColors disabled colors rendering in color function.
func disableColors() {
	colorsEnabled = false
}

// color adds given ANSII escape code at beginning of every line.
func color(str string, color colorCode) string {
	if !colorsEnabled {
		return str
	}

	return string(color) + strings.ReplaceAll(str, "\n", "\n"+string(reset+color)) + string(reset)
}
