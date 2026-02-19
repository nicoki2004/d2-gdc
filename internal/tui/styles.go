package tui

import "github.com/charmbracelet/lipgloss"

// Lipgloss styles for the TUI
var (
	// Colors
	ColorPrimary = lipgloss.Color("212")
	ColorSuccess = lipgloss.Color("42")
	ColorWarning = lipgloss.Color("208")
	ColorError   = lipgloss.Color("196")
	ColorText    = lipgloss.Color("7")
	ColorDim     = lipgloss.Color("8")
	ColorYellow  = lipgloss.Color("226")
	ColorBlue    = lipgloss.Color("4")
	ColorPurple  = lipgloss.Color("62")
	ColorGreen   = lipgloss.Color("82")

	// Base styles
	BaseStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Background(lipgloss.Color("0"))

	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorYellow).
			Background(ColorBlue).
			Padding(0, 1).
			MarginBottom(1)

	SubHeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorYellow)

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(ColorPrimary).
			Padding(0, 1)

	NormalStyle = lipgloss.NewStyle().
			Foreground(ColorText).
			Padding(0, 1)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPurple).
			Padding(1)

	BorderBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorPurple).
			Padding(0, 1)

	StatBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorSuccess).
			Padding(0, 1)

	PerkStyleGood = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	PerkStyleAvailable = lipgloss.NewStyle().
				Foreground(ColorWarning)

	PerkStyleUnavailable = lipgloss.NewStyle().
				Foreground(ColorDim).
				Strikethrough(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	DimStyle = lipgloss.NewStyle().
			Foreground(ColorDim)

	HighlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorYellow)

	FooterStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			MarginTop(1)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorDim).
			Italic(true)

	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorYellow).
				Background(ColorBlue).
				Padding(0, 1).
				Align(lipgloss.Center)

	TableRowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	TableRowSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("0")).
				Background(ColorPrimary).
				Padding(0, 1)
)

// Helper functions for centering and styling text
func CenterText(text string, width int) string {
	return lipgloss.Place(width, 1, lipgloss.Center, lipgloss.Center, text)
}

func RenderBox(content string, width int) string {
	return BorderStyle.Width(width).Render(content)
}

func RenderStatBox(label string, value string) string {
	return StatBoxStyle.Render(label + ": " + value)
}

func JoinHorizontal(items ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Left, items...)
}

func JoinVertical(items ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}
