package dispatcher

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	term       = termenv.EnvColorProfile()
	failed     = makeFgStyle("#ff0000")
	passed     = makeFgStyle("#008000")
	title      = makeFgStyle("#ffffff")
	docStyle   = lipgloss.NewStyle().Margin(1, 2)
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
	errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Render
)

// Return a function that will colorize the foreground of a given string.
func makeFgStyle(color string) func(string) string {
	return termenv.Style{}.Foreground(term.Color(color)).Styled
}
