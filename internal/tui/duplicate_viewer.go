package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DuplicateView displays duplicate weapons grouped by hash
type DuplicateView struct {
	duplicateGroups [][]WeaponData
	selectedGroup   int
	expandedGroup   int // -1 means no group is expanded
	pageSize        int
	currentPage     int
	totalPages      int
}

// NewDuplicateView creates a new duplicate viewer
func NewDuplicateView(duplicateGroups [][]WeaponData) *DuplicateView {
	d := &DuplicateView{
		duplicateGroups: duplicateGroups,
		selectedGroup:   0,
		expandedGroup:   -1,
		pageSize:        8,
		currentPage:     0,
		totalPages:      1,
	}
	d.calculatePages()
	return d
}

// SetGroups refreshes duplicate groups and resets pagination state.
func (d *DuplicateView) SetGroups(duplicateGroups [][]WeaponData) {
	d.duplicateGroups = duplicateGroups
	d.selectedGroup = 0
	d.expandedGroup = -1
	d.currentPage = 0
	d.calculatePages()
}

func (d *DuplicateView) calculatePages() {
	d.totalPages = (len(d.duplicateGroups) + d.pageSize - 1) / d.pageSize
	if d.totalPages == 0 {
		d.totalPages = 1
	}
}

// View renders the duplicate weapons view
func (d *DuplicateView) View() string {
	if len(d.duplicateGroups) == 0 {
		return BorderStyle.Render("No duplicate weapons found")
	}

	title := HeaderStyle.Render(fmt.Sprintf("📋  DUPLICATE WEAPONS (%d groups)", len(d.duplicateGroups)))

	var items []string
	start := d.currentPage * d.pageSize
	end := start + d.pageSize
	if end > len(d.duplicateGroups) {
		end = len(d.duplicateGroups)
	}

	for i := start; i < end; i++ {
		group := d.duplicateGroups[i]
		if len(group) == 0 {
			continue
		}

		marker := " "
		if i == d.selectedGroup {
			marker = "▶"
		}

		// Group header with count
		headerText := fmt.Sprintf("%s %-30s (%d copies)", marker, truncateText(group[0].Name, 30), len(group))
		header := WarningStyle.Render(headerText)
		items = append(items, header)

		// Show details if expanded
		if i == d.expandedGroup {
			for idx, weapon := range group {
				tier := d.getTierColor(weapon.Tier).Render(fmt.Sprintf("[%s]", weapon.Tier))
				power := NormalStyle.Render(fmt.Sprintf("Power:%4d", weapon.Power))
				slot := weapon.Slot
				if strings.TrimSpace(slot) == "" {
					slot = "(none)"
				}
				details := fmt.Sprintf("  ├─ [#%d] %s | %-10s | %-22s | Kills:%-5d | %s",
					idx+1, power, slot, truncateText(weapon.Location, 22), weapon.Kills, tier)
				items = append(items, details)
			}
			items = append(items, d.renderStatsComparison(group)...)
			items = append(items, "  └─ [ENTER to collapse]\n")
		} else {
			items = append(items, DimStyle.Render("  (Press ENTER to expand)\n"))
		}
	}

	content := BorderBoxStyle.Render(strings.Join(items, "\n"))
	pagination := DimStyle.Render(
		fmt.Sprintf("Page %d/%d | Showing %d-%d of %d",
			d.currentPage+1, d.totalPages, start+1, end, len(d.duplicateGroups)),
	)
	help := DimStyle.Render("[↑/↓] Navigate  [←/→] Pages  [Enter] Expand/Collapse  [ESC] Menu")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		pagination,
		"",
		help,
	)
}

// getTierColor returns the appropriate lipgloss style for a tier
func (d *DuplicateView) getTierColor(tier string) lipgloss.Style {
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

// Update handles key presses for duplicate view
func (d *DuplicateView) Update(msg tea.Msg) error {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if d.selectedGroup > 0 {
				d.selectedGroup--
				if d.selectedGroup < d.currentPage*d.pageSize {
					d.currentPage--
				}
			}
		case "down":
			if d.selectedGroup < len(d.duplicateGroups)-1 {
				d.selectedGroup++
				if d.selectedGroup >= (d.currentPage+1)*d.pageSize {
					d.currentPage++
				}
			}
		case "left":
			if d.currentPage > 0 {
				d.currentPage--
				d.selectedGroup = d.currentPage * d.pageSize
				d.expandedGroup = -1
			}
		case "right":
			if d.currentPage < d.totalPages-1 {
				d.currentPage++
				d.selectedGroup = d.currentPage * d.pageSize
				if d.selectedGroup >= len(d.duplicateGroups) {
					d.selectedGroup = len(d.duplicateGroups) - 1
				}
				d.expandedGroup = -1
			}
		case "enter":
			if d.expandedGroup == d.selectedGroup {
				d.expandedGroup = -1
			} else {
				d.expandedGroup = d.selectedGroup
			}
		}
	}
	return nil
}

// Init initializes the view
func (d *DuplicateView) Init() tea.Cmd {
	return nil
}

func (d *DuplicateView) renderStatsComparison(group []WeaponData) []string {
	commonStats := []string{
		"Range",
		"Stability",
		"Handling",
		"Reload Speed",
		"Aim Assistance",
		"Magazine",
	}

	var lines []string
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorYellow)
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorText)
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(ColorBlue)
	valueStyle := lipgloss.NewStyle().Foreground(ColorText)
	emptyStyle := lipgloss.NewStyle().Foreground(ColorDim)

	lines = append(lines, titleStyle.Render("  Stats Comparison"))

	header := fmt.Sprintf("    %-18s", "Stat")
	for i := range group {
		header += fmt.Sprintf(" %5s", fmt.Sprintf("#%d", i+1))
	}
	lines = append(lines, headerStyle.Render(header))
	lines = append(lines, lipgloss.NewStyle().Foreground(ColorPurple).Render("    "+strings.Repeat("─", len(header)-4)))

	for _, stat := range commonStats {
		var values []int
		hasAny := false
		maxValue := -1
		for _, weapon := range group {
			value := 0
			if weapon.Stats != nil {
				value = weapon.Stats[stat]
			}
			if value > 0 {
				hasAny = true
			}
			if value > maxValue {
				maxValue = value
			}
			values = append(values, value)
		}
		if !hasAny {
			continue
		}

		row := labelStyle.Render(fmt.Sprintf("    %-18s", stat))
		for _, value := range values {
			cell := fmt.Sprintf("%5s", "--")
			if value > 0 {
				cell = fmt.Sprintf("%5d", value)
			}
			if value == maxValue && maxValue > 0 {
				row += " " + PerkStyleGood.Render(cell)
			} else if value <= 0 {
				row += " " + emptyStyle.Render(cell)
			} else {
				row += " " + valueStyle.Render(cell)
			}
		}
		lines = append(lines, row)
	}

	if len(lines) == 2 {
		lines = append(lines, DimStyle.Render("    No stats found for comparison."))
	}

	return lines
}
