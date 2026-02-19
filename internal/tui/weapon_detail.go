package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WeaponDetailView displays detailed information about a single weapon
type WeaponDetailView struct {
	weapon  WeaponData
	page    int // For scrolling through large stat lists
	iconArt string
}

// NewWeaponDetailView creates a new weapon detail view
func NewWeaponDetailView(weapon WeaponData) *WeaponDetailView {
	view := &WeaponDetailView{
		weapon: weapon,
		page:   0,
	}
	view.refreshIconArt()
	return view
}

func (w *WeaponDetailView) SetWeapon(weapon WeaponData) {
	w.weapon = weapon
	w.page = 0
	w.refreshIconArt()
}

func (w *WeaponDetailView) refreshIconArt() {
	w.iconArt = ""
	if strings.TrimSpace(w.weapon.IconUrl) == "" {
		return
	}
	art, err := renderIconUnicodeFromURL(w.weapon.IconUrl, 24, 10)
	if err != nil {
		return
	}
	w.iconArt = art
}

// View renders the weapon detail with improved layout
func (w *WeaponDetailView) View() string {
	weapon := w.weapon

	// Header with tier-specific styling
	header := w.renderHeader(weapon)

	// Main information in a cohesive section
	mainInfo := w.renderMainInfo(weapon)

	// Stats with visual indicators
	statsSection := w.renderStatsSection(weapon.Stats)

	// Perks organized by socket
	perksSection := w.renderPerksSection(weapon.Perks)

	// Metadata and additional info
	metaSection := w.renderMetadata(weapon)

	// Context-specific help
	help := w.renderDetailHelp()

	renderIcon := w.renderIconPanel(weapon.IconUrl)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		mainInfo,
		lipgloss.JoinHorizontal(0.5, statsSection, renderIcon),
		perksSection,
		metaSection,
		"",
		help,
	)
}

// renderHeader creates the title with tier colors and rating
func (w *WeaponDetailView) renderHeader(weapon WeaponData) string {
	tierStyle := w.getTierColor(weapon.Tier)
	tierBadge := tierStyle.Render(fmt.Sprintf(" %s ", weapon.Tier))

	var powerIcon string
	switch {
	case weapon.Power >= 850:
		powerIcon = "⭐⭐⭐"
	case weapon.Power >= 800:
		powerIcon = "⭐⭐"
	case weapon.Power >= 750:
		powerIcon = "⭐"
	default:
		powerIcon = ""
	}

	title := fmt.Sprintf("⚔️  %s %s | 📊 Power: %d %s", weapon.Name, tierBadge, weapon.Power, powerIcon)
	return HeaderStyle.Render(title)
}

// renderMainInfo displays weapon summary in a grid
func (w *WeaponDetailView) renderMainInfo(weapon WeaponData) string {
	slot := weapon.Slot
	if strings.TrimSpace(slot) == "" {
		slot = "—"
	}

	dmgIcon := damageTypeIcon(weapon.DamageType)

	mainGrid := fmt.Sprintf(`🎯 Type:       %-20s  📍 Slot:       %s
🔥 Damage:     %s %-18s  🎮 Location:   %s
💀 Kills:      %-20d  📊 Level:      %d`,
		weapon.Type,
		slot,
		dmgIcon,
		weapon.DamageType,
		weapon.Location,
		weapon.Kills,
		weapon.Level,
	)

	return BorderBoxStyle.Render(mainGrid)
}

// renderStatsSection displays weapon stats with colored bars
func (w *WeaponDetailView) renderStatsSection(stats map[string]int) string {
	if len(stats) == 0 {
		return BorderBoxStyle.Render(SubHeaderStyle.Render("📊 Stats") + "\n" + DimStyle.Render("No stats available"))
	}

	names := make([]string, 0, len(stats))
	for name := range stats {
		names = append(names, name)
	}
	sort.Strings(names)

	lines := []string{SubHeaderStyle.Render("📊 Stats")}
	for _, name := range names {
		value := stats[name]
		bar := w.renderStatBar(value)
		colorStyle := w.getStatColor(value)
		lines = append(lines, fmt.Sprintf("  %-20s %s %s",
			name,
			colorStyle.Render(fmt.Sprintf("%3d", value)),
			bar,
		))
	}
	return BorderBoxStyle.Render(strings.Join(lines, "\n"))
}

// renderStatBar creates a visual representation of stat value
func (w *WeaponDetailView) renderStatBar(value int) string {
	bar := value / 8
	if bar > 15 {
		bar = 15
	}
	if bar < 0 {
		bar = 0
	}
	filled := strings.Repeat("█", bar)
	empty := strings.Repeat("░", 15-bar)
	return fmt.Sprintf("[%s%s]", filled, empty)
}

// getStatColor returns color based on stat value
func (w *WeaponDetailView) getStatColor(value int) lipgloss.Style {
	switch {
	case value >= 80:
		return lipgloss.NewStyle().Foreground(ColorYellow).Bold(true)
	case value >= 60:
		return lipgloss.NewStyle().Foreground(ColorGreen)
	case value >= 40:
		return lipgloss.NewStyle().Foreground(ColorBlue)
	default:
		return lipgloss.NewStyle().Foreground(ColorDim)
	}
}

// renderPerksSection displays perks organized by socket
func (w *WeaponDetailView) renderPerksSection(perks []PerkData) string {
	if len(perks) == 0 {
		return BorderBoxStyle.Render(SubHeaderStyle.Render("⚡ Perks") + "\n" + DimStyle.Render("No perks available"))
	}

	sort.Slice(perks, func(i, j int) bool {
		if perks[i].SocketIdx == perks[j].SocketIdx {
			if perks[i].IsEquipped == perks[j].IsEquipped {
				return perks[i].Name < perks[j].Name
			}
			return perks[i].IsEquipped
		}
		return perks[i].SocketIdx < perks[j].SocketIdx
	})

	lines := []string{SubHeaderStyle.Render("⚡ Perks")}
	currentSocket := -1
	for _, perk := range perks {
		if perk.SocketIdx != currentSocket {
			currentSocket = perk.SocketIdx
			lines = append(lines, HighlightStyle.Render(fmt.Sprintf("  🔌 Socket %d", currentSocket+1)))
		}

		marker := "○"
		style := PerkStyleAvailable
		if perk.IsEquipped {
			marker = "●"
			style = PerkStyleGood
		}
		perkName := truncateText(perk.Name, 70)
		lines = append(lines, "    "+style.Render(fmt.Sprintf("%s %s", marker, perkName)))
	}
	return BorderBoxStyle.Render(strings.Join(lines, "\n"))
}

// renderMetadata displays additional weapon metadata
func (w *WeaponDetailView) renderMetadata(weapon WeaponData) string {
	instance := weapon.InstanceID
	if len(instance) > 16 {
		instance = instance[:16] + "..."
	}

	iconDisplay := "URL not available"
	if weapon.IconUrl != "" {
		iconDisplay = truncateText(weapon.IconUrl, 60)
	}

	meta := fmt.Sprintf(`📋 Instance:  %s
🖼️  Icon URL:  %s`,
		instance,
		iconDisplay,
	)

	return DimStyle.Render(meta)
}

// renderDetailHelp returns context-specific help text
func (w *WeaponDetailView) renderDetailHelp() string {
	return HelpStyle.Render("↑/↓: Scroll | Enter: List | Esc: Home | ?:Help")
}

func (w *WeaponDetailView) renderIconPanel(iconURL string) string {
	if w.iconArt != "" {
		return BorderBoxStyle.Render(w.iconArt)
		// return w.iconArt
	}
	return fmt.Sprintf("Icon preview unavailable\n%s", truncateText(iconURL, 64))
}

// getTierColor returns the appropriate lipgloss style for a tier
func (w *WeaponDetailView) getTierColor(tier string) lipgloss.Style {
	switch tier {
	case "Exotic":
		return lipgloss.NewStyle().Foreground(ColorYellow).Bold(true)
	case "Legendary":
		return lipgloss.NewStyle().Foreground(ColorPurple).Bold(true)
	case "Rare":
		return lipgloss.NewStyle().Foreground(ColorBlue)
	case "Uncommon":
		return lipgloss.NewStyle().Foreground(ColorGreen)
	default:
		return lipgloss.NewStyle().Foreground(ColorText)
	}
}

// Update handles key presses for detail view
func (w *WeaponDetailView) Update(msg tea.Msg) error {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if w.page > 0 {
				w.page--
			}
		case "down":
			w.page++
		}
	}
	return nil
}

// Init initializes the view
func (w *WeaponDetailView) Init() tea.Cmd {
	return nil
}
