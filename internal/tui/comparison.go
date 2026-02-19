package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ComparisonView compares up to 4 weapons side-by-side.
type ComparisonView struct {
	weapons       []WeaponData
	allWeapons    []WeaponData
	selectedOrder []string
	selectedByID  map[string]WeaponData
	selectedIdx   int
	maxSelected   int
	pageSize      int
	currentPage   int
	totalPages    int
	searchInput   string
	isFiltering   bool
	filterApplied bool
	rarityOptions []string
	rarityIdx     int
	typeOptions   []string
	typeIdx       int
	damageOptions []string
	damageIdx     int
}

// NewComparisonView creates a new comparison view.
func NewComparisonView(maxSelected int) *ComparisonView {
	if maxSelected <= 0 {
		maxSelected = 4
	}
	if maxSelected > 4 {
		maxSelected = 4
	}
	return &ComparisonView{
		weapons:       []WeaponData{},
		allWeapons:    []WeaponData{},
		selectedOrder: []string{},
		selectedByID:  map[string]WeaponData{},
		selectedIdx:   0,
		maxSelected:   maxSelected,
		pageSize:      10,
		currentPage:   0,
		totalPages:    1,
		searchInput:   "",
		isFiltering:   false,
		filterApplied: false,
		rarityOptions: []string{"All", "Exotic", "Legendary"},
		rarityIdx:     0,
		typeOptions:   []string{"All"},
		typeIdx:       0,
		damageOptions: []string{"All", "Normal", "Solar", "Arc", "Void", "Strand", "Stasis"},
		damageIdx:     0,
	}
}

// SetWeapons updates the pool of weapons shown in the selector.
func (c *ComparisonView) SetWeapons(weapons []WeaponData) {
	c.allWeapons = append([]WeaponData(nil), weapons...)
	c.weapons = append([]WeaponData(nil), weapons...)
	sort.Slice(c.allWeapons, func(i, j int) bool {
		if c.allWeapons[i].Power == c.allWeapons[j].Power {
			return c.allWeapons[i].Name < c.allWeapons[j].Name
		}
		return c.allWeapons[i].Power > c.allWeapons[j].Power
	})
	c.weapons = append([]WeaponData(nil), c.allWeapons...)

	// Keep selected weapons only if still present.
	indexByID := make(map[string]WeaponData, len(c.allWeapons))
	for _, w := range c.allWeapons {
		indexByID[w.InstanceID] = w
	}
	newOrder := make([]string, 0, len(c.selectedOrder))
	for _, id := range c.selectedOrder {
		if _, ok := indexByID[id]; ok {
			newOrder = append(newOrder, id)
		} else {
			delete(c.selectedByID, id)
		}
	}
	c.selectedOrder = newOrder

	c.searchInput = ""
	c.isFiltering = false
	c.filterApplied = false
	c.rarityIdx = 0
	c.typeIdx = 0
	c.damageIdx = 0
	c.refreshTypeOptions()

	if c.selectedIdx >= len(c.weapons) {
		c.selectedIdx = max(0, len(c.weapons)-1)
	}
	c.calculatePages()
	c.alignCurrentPage()
}

func (c *ComparisonView) calculatePages() {
	c.totalPages = (len(c.weapons) + c.pageSize - 1) / c.pageSize
	if c.totalPages == 0 {
		c.totalPages = 1
	}
}

func (c *ComparisonView) alignCurrentPage() {
	if c.selectedIdx < c.currentPage*c.pageSize {
		c.currentPage = c.selectedIdx / c.pageSize
	}
	for c.selectedIdx >= (c.currentPage+1)*c.pageSize && c.currentPage < c.totalPages-1 {
		c.currentPage++
	}
}

// CurrentWeapon returns the cursor weapon.
func (c *ComparisonView) CurrentWeapon() (WeaponData, bool) {
	if c.selectedIdx < 0 || c.selectedIdx >= len(c.weapons) {
		return WeaponData{}, false
	}
	return c.weapons[c.selectedIdx], true
}

// IsSelected tells if a weapon is currently in the comparison set.
func (c *ComparisonView) IsSelected(instanceID string) bool {
	_, ok := c.selectedByID[instanceID]
	return ok
}

// UpsertSelectedWeapon replaces or inserts the selected weapon with enriched data.
func (c *ComparisonView) UpsertSelectedWeapon(weapon WeaponData) {
	if !c.IsSelected(weapon.InstanceID) {
		return
	}
	c.selectedByID[weapon.InstanceID] = weapon
}

func (c *ComparisonView) selectedWeapons() []WeaponData {
	result := make([]WeaponData, 0, len(c.selectedOrder))
	for _, id := range c.selectedOrder {
		if w, ok := c.selectedByID[id]; ok {
			result = append(result, w)
		}
	}
	return result
}

func (c *ComparisonView) toggleCurrentSelection() {
	weapon, ok := c.CurrentWeapon()
	if !ok {
		return
	}

	if c.IsSelected(weapon.InstanceID) {
		delete(c.selectedByID, weapon.InstanceID)
		newOrder := make([]string, 0, len(c.selectedOrder))
		for _, id := range c.selectedOrder {
			if id != weapon.InstanceID {
				newOrder = append(newOrder, id)
			}
		}
		c.selectedOrder = newOrder
		return
	}

	if len(c.selectedOrder) >= c.maxSelected {
		return
	}
	c.selectedOrder = append(c.selectedOrder, weapon.InstanceID)
	c.selectedByID[weapon.InstanceID] = weapon
}

// View renders the selector and comparison output.
func (c *ComparisonView) View() string {
	title := HeaderStyle.Render(fmt.Sprintf("⚖️  WEAPON COMPARISON (%d/%d selected)", len(c.selectedOrder), c.maxSelected))

	if len(c.weapons) == 0 {
		filterInfo := HighlightStyle.Render(fmt.Sprintf("Rarity[r]: %s | Type[t]: %s | Damage[d]: %s | Search: %q",
			c.currentRarityFilter(), c.currentTypeFilter(), c.currentDamageFilter(), c.searchInput))
		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			BorderStyle.Render("No weapons found for current filters"),
			filterInfo,
			DimStyle.Render("[/] Search  [r] Rarity  [t] Type  [d] Damage  [Backspace] Clear all"),
		)
	}

	start := c.currentPage * c.pageSize
	end := start + c.pageSize
	if end > len(c.weapons) {
		end = len(c.weapons)
	}

	rows := []string{
		SubHeaderStyle.Render("Sel Cur  Name                          Slot       Power Tier       Dmg"),
		DimStyle.Render(strings.Repeat("─", 78)),
	}

	for i := start; i < end; i++ {
		w := c.weapons[i]
		sel := "[ ]"
		if c.IsSelected(w.InstanceID) {
			sel = "[X]"
		}
		cur := " "
		if i == c.selectedIdx {
			cur = "▶"
		}
		slot := w.Slot
		if strings.TrimSpace(slot) == "" {
			slot = "(none)"
		}
		dmgIcon := damageTypeIcon(w.DamageType)
		tierText := formatTierCell(w.Tier)
		tierStr := c.getTierColor(w.Tier).Render(tierText)
		line := fmt.Sprintf("%s  %s   %-28s %-10s %5d %s %s",
			sel,
			cur,
			truncateText(w.Name, 28),
			truncateText(slot, 10),
			w.Power,
			tierStr,
			dmgIcon,
		)
		if i == c.selectedIdx {
			rows = append(rows, SelectedStyle.Render(line))
		} else {
			rows = append(rows, NormalStyle.Render(line))
		}
	}

	selectorPanel := BorderBoxStyle.Render(strings.Join(rows, "\n"))
	pagination := DimStyle.Render(fmt.Sprintf("Page %d/%d | Showing %d-%d of %d",
		c.currentPage+1, c.totalPages, start+1, end, len(c.weapons)))
	quickFilters := HighlightStyle.Render(fmt.Sprintf("Rarity[r]: %s | Type[t]: %s | Damage[d]: %s",
		c.currentRarityFilter(), c.currentTypeFilter(), c.currentDamageFilter()))
	filterInfo := ""
	if c.isFiltering {
		filterInfo = HighlightStyle.Render(fmt.Sprintf("Search: %s_", c.searchInput))
	} else if c.filterApplied {
		filterInfo = DimStyle.Render(fmt.Sprintf("Search: %s", c.searchInput))
	}

	comparePanel := BorderBoxStyle.Render(c.renderComparisonDetails())

	help := DimStyle.Render("[↑/↓] Move  [←/→] Pages  [Space/Enter] Select  [/] Search  [r] Rarity  [t] Type  [d] Damage  [Backspace] Remove last")
	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		selectorPanel,
		pagination,
		quickFilters,
		filterInfo,
		"",
		comparePanel,
		"",
		help,
	)
}

func (c *ComparisonView) renderComparisonDetails() string {
	selected := c.selectedWeapons()
	if len(selected) == 0 {
		return fmt.Sprintf("Select up to %d weapons to compare stats and perks.", c.maxSelected)
	}

	lines := []string{
		SubHeaderStyle.Render("Stats Comparison"),
	}
	lines = append(lines, c.renderStatsTable(selected)...)
	lines = append(lines, "")
	lines = append(lines, SubHeaderStyle.Render("Perks Comparison (equipped)"))
	lines = append(lines, c.renderPerksTable(selected)...)
	return strings.Join(lines, "\n")
}

func (c *ComparisonView) renderStatsTable(weapons []WeaponData) []string {
	priority := []string{
		"Impact", "Range", "Stability", "Handling", "Reload Speed", "Aim Assistance", "Magazine", "Zoom", "RPM",
	}
	statsSet := map[string]bool{}
	for _, w := range weapons {
		for name, value := range w.Stats {
			if value > 0 {
				statsSet[name] = true
			}
		}
	}

	var statNames []string
	for _, name := range priority {
		if statsSet[name] {
			statNames = append(statNames, name)
			delete(statsSet, name)
		}
	}
	var rest []string
	for name := range statsSet {
		rest = append(rest, name)
	}
	sort.Strings(rest)
	statNames = append(statNames, rest...)

	if len(statNames) == 0 {
		return []string{DimStyle.Render("No stats available for selected weapons.")}
	}

	header := fmt.Sprintf("%-18s", "Stat")
	for _, w := range weapons {
		header += " " + fmt.Sprintf("%10s", truncateText(w.Name, 10))
	}
	lines := []string{HighlightStyle.Render(header)}
	lines = append(lines, DimStyle.Render(strings.Repeat("─", len(header))))

	for _, stat := range statNames {
		maxVal := -1
		values := make([]int, len(weapons))
		for i, w := range weapons {
			v := 0
			if w.Stats != nil {
				v = w.Stats[stat]
			}
			values[i] = v
			if v > maxVal {
				maxVal = v
			}
		}
		row := fmt.Sprintf("%-18s", stat)
		for _, v := range values {
			cell := fmt.Sprintf("%10s", "--")
			if v > 0 {
				cell = fmt.Sprintf("%10d", v)
			}
			if v == maxVal && maxVal > 0 {
				row += PerkStyleGood.Render(cell)
			} else {
				row += DimStyle.Render(cell)
			}
		}
		lines = append(lines, row)
	}
	return lines
}

func (c *ComparisonView) renderPerksTable(weapons []WeaponData) []string {
	perksByWeapon := make([][]string, len(weapons))
	maxRows := 0
	for i, w := range weapons {
		perksByWeapon[i] = equippedPerks(w.Perks)
		if len(perksByWeapon[i]) > maxRows {
			maxRows = len(perksByWeapon[i])
		}
	}

	if maxRows == 0 {
		return []string{DimStyle.Render("No equipped perks available for selected weapons.")}
	}

	header := fmt.Sprintf("%-8s", "Slot")
	for _, w := range weapons {
		header += " " + fmt.Sprintf("%18s", truncateText(w.Name, 18))
	}
	lines := []string{HighlightStyle.Render(header)}
	lines = append(lines, DimStyle.Render(strings.Repeat("─", len(header))))

	for i := 0; i < maxRows; i++ {
		row := fmt.Sprintf("%-8s", fmt.Sprintf("Perk %d", i+1))
		for _, perks := range perksByWeapon {
			value := "--"
			if i < len(perks) {
				value = perks[i]
			}
			row += " " + fmt.Sprintf("%18s", truncateText(value, 18))
		}
		lines = append(lines, row)
	}

	return lines
}

func equippedPerks(perks []PerkData) []string {
	if len(perks) == 0 {
		return nil
	}

	filtered := make([]PerkData, 0, len(perks))
	for _, p := range perks {
		if p.IsEquipped {
			filtered = append(filtered, p)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].SocketIdx == filtered[j].SocketIdx {
			return filtered[i].Name < filtered[j].Name
		}
		return filtered[i].SocketIdx < filtered[j].SocketIdx
	})

	result := make([]string, 0, len(filtered))
	for _, p := range filtered {
		result = append(result, p.Name)
	}
	return result
}

func (c *ComparisonView) getTierColor(tier string) lipgloss.Style {
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

// Update handles key presses for comparison selector.
func (c *ComparisonView) Update(msg tea.Msg) error {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.isFiltering {
			switch msg.String() {
			case "enter":
				c.isFiltering = false
			case "esc":
				c.isFiltering = false
			case "backspace":
				if len(c.searchInput) > 0 {
					c.searchInput = c.searchInput[:len(c.searchInput)-1]
				}
				c.applyFilter()
			default:
				if len(msg.Runes) > 0 {
					c.searchInput += string(msg.Runes)
					c.applyFilter()
				}
			}
			return nil
		}

		switch msg.String() {
		case "/":
			c.isFiltering = true
		case "r":
			c.rarityIdx = (c.rarityIdx + 1) % len(c.rarityOptions)
			c.applyFilter()
		case "t":
			c.typeIdx = (c.typeIdx + 1) % len(c.typeOptions)
			c.applyFilter()
		case "d":
			c.damageIdx = (c.damageIdx + 1) % len(c.damageOptions)
			c.applyFilter()
		case "up":
			if c.selectedIdx > 0 {
				c.selectedIdx--
				c.alignCurrentPage()
			}
		case "down":
			if c.selectedIdx < len(c.weapons)-1 {
				c.selectedIdx++
				c.alignCurrentPage()
			}
		case "left":
			if c.currentPage > 0 {
				c.currentPage--
				c.selectedIdx = c.currentPage * c.pageSize
			}
		case "right":
			if c.currentPage < c.totalPages-1 {
				c.currentPage++
				c.selectedIdx = c.currentPage * c.pageSize
				if c.selectedIdx >= len(c.weapons) {
					c.selectedIdx = len(c.weapons) - 1
				}
			}
		case "enter", " ":
			c.toggleCurrentSelection()
		case "backspace":
			if len(c.selectedOrder) > 0 {
				lastID := c.selectedOrder[len(c.selectedOrder)-1]
				delete(c.selectedByID, lastID)
				c.selectedOrder = c.selectedOrder[:len(c.selectedOrder)-1]
			} else if c.filterApplied {
				c.clearFilter()
			}
		}
	}
	return nil
}

func (c *ComparisonView) IsFiltering() bool {
	return c.isFiltering
}

func (c *ComparisonView) applyFilter() {
	query := strings.TrimSpace(strings.ToLower(c.searchInput))
	rarityFilter := strings.ToLower(c.currentRarityFilter())
	typeFilter := strings.ToLower(c.currentTypeFilter())
	damageFilter := strings.ToLower(c.currentDamageFilter())

	filtered := make([]WeaponData, 0, len(c.allWeapons))
	for _, weapon := range c.allWeapons {
		if rarityFilter != "all" && strings.ToLower(strings.TrimSpace(weapon.Tier)) != rarityFilter {
			continue
		}
		if typeFilter != "all" && strings.ToLower(strings.TrimSpace(weapon.Type)) != typeFilter {
			continue
		}
		if damageFilter != "all" && strings.ToLower(normalizedDamageType(weapon.DamageType)) != damageFilter {
			continue
		}
		if query != "" {
			haystack := strings.ToLower(weapon.Name + " " + weapon.Type + " " + weapon.Tier + " " + weapon.Slot)
			if !strings.Contains(haystack, query) {
				continue
			}
		}
		filtered = append(filtered, weapon)
	}

	c.weapons = filtered
	c.filterApplied = c.hasActiveFilter()
	c.selectedIdx = 0
	c.currentPage = 0
	c.calculatePages()
}

func (c *ComparisonView) clearFilter() {
	c.weapons = append([]WeaponData(nil), c.allWeapons...)
	c.searchInput = ""
	c.rarityIdx = 0
	c.typeIdx = 0
	c.damageIdx = 0
	c.filterApplied = false
	c.selectedIdx = 0
	c.currentPage = 0
	c.calculatePages()
}

func (c *ComparisonView) refreshTypeOptions() {
	typeSet := map[string]struct{}{}
	for _, weapon := range c.allWeapons {
		t := strings.TrimSpace(weapon.Type)
		if t == "" {
			continue
		}
		typeSet[t] = struct{}{}
	}

	types := make([]string, 0, len(typeSet))
	for t := range typeSet {
		types = append(types, t)
	}
	sort.Strings(types)

	c.typeOptions = append([]string{"All"}, types...)
	if c.typeIdx >= len(c.typeOptions) {
		c.typeIdx = 0
	}
}

func (c *ComparisonView) currentRarityFilter() string {
	if c.rarityIdx < 0 || c.rarityIdx >= len(c.rarityOptions) {
		return "All"
	}
	return c.rarityOptions[c.rarityIdx]
}

func (c *ComparisonView) currentTypeFilter() string {
	if c.typeIdx < 0 || c.typeIdx >= len(c.typeOptions) {
		return "All"
	}
	return c.typeOptions[c.typeIdx]
}

func (c *ComparisonView) hasActiveFilter() bool {
	return strings.TrimSpace(c.searchInput) != "" ||
		!strings.EqualFold(c.currentRarityFilter(), "All") ||
		!strings.EqualFold(c.currentTypeFilter(), "All") ||
		!strings.EqualFold(c.currentDamageFilter(), "All")
}

func (c *ComparisonView) currentDamageFilter() string {
	if c.damageIdx < 0 || c.damageIdx >= len(c.damageOptions) {
		return "All"
	}
	return c.damageOptions[c.damageIdx]
}

// Init initializes the view.
func (c *ComparisonView) Init() tea.Cmd {
	return nil
}
