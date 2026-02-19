package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GodRollView searches and displays god rolls
type GodRollView struct {
	godRollDefinitions map[string][]string // Weapon name -> ideal perks
	results            []WeaponData
	selectedResult     int
}

// GodRollDefinition represents ideal perks for a weapon
type GodRollDefinition struct {
	WeaponName    string
	DesiredPerks  []string
	MinStatValues map[string]int
}

// NewGodRollView creates a new god roll search view
func NewGodRollView() *GodRollView {
	return &GodRollView{
		godRollDefinitions: initializeGodRolls(),
		results:            make([]WeaponData, 0),
		selectedResult:     0,
	}
}

// initializeGodRolls sets up common god roll definitions
func initializeGodRolls() map[string][]string {
	return map[string][]string{
		"Witherhoard":      {"Volatile Rounds", "Frenzy"},
		"Tarrabah":         {"Smooth Ballistics", "Kill Clip"},
		"Legend of Acrius": {"Precision Slug", "Rampage"},
		"Duality":          {"Void", "One-Two Punch"},
		"Commemoration":    {"Subsistence", "Dragonfly"},
	}
}

// Search finds weapons matching god roll criteria
func (g *GodRollView) Search(allWeapons []WeaponData, weaponName string) {
	g.results = make([]WeaponData, 0)
	g.selectedResult = 0

	idealPerks, exists := g.godRollDefinitions[weaponName]
	if !exists {
		return
	}

	for _, weapon := range allWeapons {
		if weapon.Name == weaponName {
			score := g.calculateScore(weapon, idealPerks)
			if score >= 75 { // Only show results >= 75%
				g.results = append(g.results, weapon)
			}
		}
	}
}

// calculateScore calculates how well a weapon matches the god roll
func (g *GodRollView) calculateScore(weapon WeaponData, idealPerks []string) float64 {
	if len(idealPerks) == 0 {
		return 0
	}

	matches := 0
	for _, desiredPerk := range idealPerks {
		for _, equippedPerk := range weapon.Perks {
			if equippedPerk.Name == desiredPerk && equippedPerk.IsEquipped {
				matches++
			}
		}
	}

	return float64(matches) / float64(len(idealPerks)) * 100
}

// View renders the god roll search results
func (g *GodRollView) View() string {
	title := HeaderStyle.Render("✨  GOD ROLL SEARCH")

	if len(g.results) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			BorderStyle.Render("No god rolls found matching your criteria"),
		)
	}

	var items []string

	for i, weapon := range g.results {
		score := g.calculateScore(weapon, g.godRollDefinitions[weapon.Name])
		marker := " "
		if i == g.selectedResult {
			marker = "▶"
		}

		// Color code based on match percentage
		var scoreStyle lipgloss.Style
		switch {
		case score >= 90:
			scoreStyle = PerkStyleGood
		case score >= 75:
			scoreStyle = PerkStyleAvailable
		default:
			scoreStyle = DimStyle
		}

		scoreStr := scoreStyle.Render(fmt.Sprintf("%.0f%%", score))
		line := fmt.Sprintf("%s [%s Match] %s [Power: %d]",
			marker, scoreStr, weapon.Name, weapon.Power)

		if i == g.selectedResult {
			items = append(items, SelectedStyle.Render(line))
		} else {
			items = append(items, NormalStyle.Render(line))
		}
	}

	content := BorderBoxStyle.Render(strings.Join(items, "\n"))
	help := DimStyle.Render("[↑/↓] Navigate  [Enter] Details  [ESC] Menu")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		"",
		help,
	)
}

// Update handles key presses for god roll view
func (g *GodRollView) Update(msg tea.Msg) error {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if g.selectedResult > 0 {
				g.selectedResult--
			}
		case "down":
			if g.selectedResult < len(g.results)-1 {
				g.selectedResult++
			}
		}
	}
	return nil
}

// Init initializes the view
func (g *GodRollView) Init() tea.Cmd {
	return nil
}
