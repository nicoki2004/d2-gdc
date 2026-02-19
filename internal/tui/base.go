package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// View defines the interface all views must implement
type View interface {
	View() string
	Update(msg tea.Msg) error
	Init() tea.Cmd
}

// HomeView is the main menu view
type HomeView struct {
	selected        int
	totalWeapons    int
	totalCharacters int
}

// NewHomeView creates a new home view
func NewHomeView() *HomeView {
	return &HomeView{selected: 0}
}

func (h *HomeView) View() string {
	title := HeaderStyle.Render("🔫  DESTINY 2 GOD ROLL CHECKER")

	menuItems := []string{
		"Search Weapons",
		"Find Duplicates",
		"Compare Weapons",
		"Search God Rolls",
		"Help",
		"Quit",
	}

	var menu string
	menu = SubHeaderStyle.Render("Main Menu\n")
	for i, item := range menuItems {
		marker := " "
		if i == h.selected {
			marker = "▶"
		}

		key := ""
		switch i {
		case 0:
			key = "[1]"
		case 1:
			key = "[2]"
		case 2:
			key = "[3]"
		case 3:
			key = "[4]"
		case 4:
			key = "[h]"
		case 5:
			key = "[q]"
		}

		line := fmt.Sprintf("%s %s  %s", marker, key, item)
		if i == h.selected {
			menu += SelectedStyle.Render(line) + "\n"
		} else {
			menu += NormalStyle.Render(line) + "\n"
		}
	}

	stats := fmt.Sprintf("Total Weapons: %d | Characters: %d", h.totalWeapons, h.totalCharacters)
	statsBox := BorderBoxStyle.Render(stats)

	help := DimStyle.Render("[↑/↓] Navigate  [Enter/Number] Select  [ESC] Quit")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		BorderBoxStyle.Render(menu),
		"",
		statsBox,
		"",
		help,
	)
}

func (h *HomeView) Update(msg tea.Msg) error {
	return nil
}

func (h *HomeView) Init() tea.Cmd {
	return nil
}
