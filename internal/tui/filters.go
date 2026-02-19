package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func normalizedDamageType(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "solar":
		return "Solar"
	case "arc":
		return "Arc"
	case "void":
		return "Void"
	case "strand":
		return "Strand"
	case "stasis":
		return "Stasis"
	case "kinetic", "normal", "none", "":
		return "Normal"
	default:
		return "Normal"
	}
}

func damageTypeIcon(raw string) string {
	damage := normalizedDamageType(raw)
	switch damage {
	case "Solar":
		return lipgloss.NewStyle().Foreground(ColorWarning).Render("●")
	case "Arc":
		return lipgloss.NewStyle().Foreground(ColorBlue).Render("●")
	case "Void":
		return lipgloss.NewStyle().Foreground(ColorPurple).Render("●")
	case "Strand":
		return lipgloss.NewStyle().Foreground(ColorGreen).Render("●")
	case "Stasis":
		return lipgloss.NewStyle().Foreground(lipgloss.Color("123")).Render("●")
	default:
		return lipgloss.NewStyle().Foreground(ColorText).Render("●")
	}
}
