package colors

// This is taken from zerolog's repo and will be used to colorize log output
// Source: https://github.com/rs/zerolog/blob/4fff5db29c3403bc26dee9895e12a108aacc0203/console.go
const (
	// COLOR_BLACK is the ANSI code for black
	COLOR_BLACK = iota + 30
	// COLOR_RED is the ANSI code for red
	COLOR_RED
	// COLOR_GREEN is the ANSI code for green
	COLOR_GREEN
	// COLOR_YELLOW is the ANSI code for yellow
	COLOR_YELLOW
	// COLOR_BLUE is the ANSI code for blue
	COLOR_BLUE
	// COLOR_MAGENTA is the ANSI code for magenta
	COLOR_MAGENTA
	// COLOR_CYAN is the ANSI code for cyan
	COLOR_CYAN
	// COLOR_WHITE is the ANSI code for white
	COLOR_WHITE
	// COLOR_BOLD is the ANSI code for bold tet
	COLOR_BOLD = 1
	// COLOR_DARK_GRAY is the ANSI code for dark gray
	COLOR_DARK_GRAY = 90
)

// This enum is to identify special unicode characters that will be used for pretty console output
const (
	// LEFT_ARROW is the unicode string for a left arrow glyph
	LEFT_ARROW = "\u21fe"

	// DOWNWARD_LEFT_ARROW is the unicode string for a left arrow glyph but with a downward intention
	DOWNWARD_LEFT_ARROW = "\u21b3"
)
