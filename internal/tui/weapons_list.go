package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FilterOptions holds the available options for each filter type
type FilterOptions struct {
	Rarity map[string]int // "All", "Exotic", "Legendary"
	Type   map[string]int
	Damage map[string]int // "All", "Normal", "Solar", "Arc", "Void", "Strand", "Stasis"
}

// FilterState tracks current filter selections
type FilterState struct {
	RarityIdx int
	TypeIdx   int
	DamageIdx int
}

// WeaponsListView displays a list of weapons
type WeaponsListView struct {
	Weapons       []WeaponData
	SelectedIdx   int
	searchInput   string
	allWeapons    []WeaponData
	isFiltering   bool
	filterApplied bool

	// Filter state and options
	filterState   FilterState
	rarityOptions []string
	typeOptions   []string
	damageOptions []string

	// Pagination
	PageSize    int
	CurrentPage int
	TotalPages  int

	// Dimensions
	Width  int
	Height int
}

// NewWeaponsListView creates a new weapons list view
func NewWeaponsListView(weapons []WeaponData, width, height int) *WeaponsListView {
	wlv := &WeaponsListView{
		Weapons:       weapons,
		allWeapons:    append([]WeaponData(nil), weapons...),
		SelectedIdx:   0,
		PageSize:      10,
		Width:         width,
		Height:        height,
		searchInput:   "",
		isFiltering:   false,
		filterApplied: false,
		rarityOptions: []string{"All", "Exotic", "Legendary"},
		damageOptions: []string{"All", "Normal", "Solar", "Arc", "Void", "Strand", "Stasis"},
		typeOptions:   []string{"All"},
		filterState:   FilterState{RarityIdx: 0, TypeIdx: 0, DamageIdx: 0},
	}
	wlv.refreshTypeOptions()
	wlv.calculatePages()
	return wlv
}

// SetWeapons updates the weapons list
func (w *WeaponsListView) SetWeapons(weapons []WeaponData) {
	w.allWeapons = append([]WeaponData(nil), weapons...)
	w.Weapons = append([]WeaponData(nil), weapons...)
	w.SelectedIdx = 0
	w.searchInput = ""
	w.isFiltering = false
	w.filterApplied = false
	w.filterState = FilterState{RarityIdx: 0, TypeIdx: 0, DamageIdx: 0}
	w.refreshTypeOptions()
	w.CurrentPage = 0
	w.calculatePages()
}

// calculatePages calculates the number of pages
func (w *WeaponsListView) calculatePages() {
	w.TotalPages = (len(w.Weapons) + w.PageSize - 1) / w.PageSize
	if w.TotalPages == 0 {
		w.TotalPages = 1
	}
}

// getPageRange returns the start and end indices for the current page
func (w *WeaponsListView) getPageRange() (int, int) {
	start := w.CurrentPage * w.PageSize
	end := start + w.PageSize
	if end > len(w.Weapons) {
		end = len(w.Weapons)
	}
	return start, end
}

// View renders the weapons list with pagination, filtering, and navigation help
func (w *WeaponsListView) View() string {
	title := w.renderTitle()

	if len(w.Weapons) == 0 {
		return w.renderEmptyState(title)
	}

	content := w.renderWeaponsList()
	status := w.renderStatus()
	quickFilters := w.renderQuickFilters()
	help := w.renderHelpText()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		status,
		quickFilters,
		help,
	)
}

// renderTitle renders the main title with weapon count
func (w *WeaponsListView) renderTitle() string {
	var icon string
	if len(w.Weapons) > 0 {
		icon = "⚔️"
	} else {
		icon = "📋"
	}
	return HeaderStyle.Render(fmt.Sprintf("%s  WEAPONS (%d total)", icon, len(w.Weapons)))
}

// renderQuickFilters renders the active filters in a compact way
func (w *WeaponsListView) renderQuickFilters() string {
	var filters []string

	if w.currentRarityFilter() != "All" {
		filters = append(filters, fmt.Sprintf("🎖️  %s", w.currentRarityFilter()))
	}
	if w.currentTypeFilter() != "All" {
		filters = append(filters, fmt.Sprintf("🔫 %s", w.currentTypeFilter()))
	}
	if w.currentDamageFilter() != "All" {
		filters = append(filters, fmt.Sprintf("⚡ %s", w.currentDamageFilter()))
	}

	if len(filters) == 0 {
		return ""
	}

	filterText := strings.Join(filters, " | ")
	return HighlightStyle.Render("📍 Active: " + filterText)
}

// renderStatus renders pagination and filter status
func (w *WeaponsListView) renderStatus() string {
	start, end := w.getPageRange()
	pageInfo := fmt.Sprintf("📄 Page %d/%d | Items %d-%d/%d", w.CurrentPage+1, w.TotalPages, start+1, end, len(w.Weapons))

	var searchStatus string
	if w.isFiltering {
		searchStatus = fmt.Sprintf(" | 🔍 Searching: %s_", w.searchInput)
	} else if w.filterApplied {
		searchStatus = fmt.Sprintf(" | 🔍 Filter: \"%s\"", w.searchInput)
	}

	return DimStyle.Render(pageInfo + searchStatus)
}

// renderHelpText renders keyboard shortcuts in a table format
func (w *WeaponsListView) renderHelpText() string {
	var shortcuts []string

	if w.isFiltering {
		shortcuts = []string{
			"Enter: Done",
			"Esc: Cancel",
			"Bkspc: Delete",
		}
	} else {
		shortcuts = []string{
			"↑/↓: Select | ←/→: Pages | /: Search",
			"r: Rarity | t: Type | d: Damage | Bkspc: Clear",
			"Enter: Details | Esc: Menu",
		}
	}

	return HelpStyle.Render(strings.Join(shortcuts, "  |  "))
}

// renderEmptyState renders the view when no weapons match filters
func (w *WeaponsListView) renderEmptyState(title string) string {
	message := BorderStyle.Render("📭 No weapons found for current filters")

	var activeFilters []string
	if w.currentRarityFilter() != "All" {
		activeFilters = append(activeFilters, fmt.Sprintf("Rarity: %s", w.currentRarityFilter()))
	}
	if w.currentTypeFilter() != "All" {
		activeFilters = append(activeFilters, fmt.Sprintf("Type: %s", w.currentTypeFilter()))
	}
	if w.currentDamageFilter() != "All" {
		activeFilters = append(activeFilters, fmt.Sprintf("Damage: %s", w.currentDamageFilter()))
	}
	if strings.TrimSpace(w.searchInput) != "" {
		activeFilters = append(activeFilters, fmt.Sprintf("Search: \"%s\"", w.searchInput))
	}

	var hint string
	if len(activeFilters) > 0 {
		hint = HighlightStyle.Render("Active: " + strings.Join(activeFilters, " | "))
	} else {
		hint = DimStyle.Render("Try adding filters or search terms")
	}

	help := HelpStyle.Render("[Bkspc] Clear filters | [/] Search | [Esc] Menu")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		message,
		hint,
		"",
		help,
	)
}

// renderWeaponsList renders a responsive weapons table
func (w *WeaponsListView) renderWeaponsList() string {
	head := w.buildTableHeader()
	separator := DimStyle.Render(strings.Repeat("─", w.getEffectiveWidth()))

	var items []string
	start, end := w.getPageRange()
	for i := start; i < end; i++ {
		weapon := w.Weapons[i]
		line := w.formatWeaponLine(weapon, i == w.SelectedIdx)
		items = append(items, line)
	}

	itemsText := strings.Join(items, "\n")
	content := lipgloss.JoinVertical(lipgloss.Left, head, separator, itemsText)

	// Add scroll indicator
	if w.hasMoreItems() {
		content += "\n" + w.renderScrollIndicator()
	}

	return BorderBoxStyle.Width(w.getEffectiveWidth()).Render(content)
}

// buildTableHeader constructs the header row with dynamic columns
func (w *WeaponsListView) buildTableHeader() string {
	header := fmt.Sprintf("%-3s  %-28s  %-9s  %-11s %s",
		"Sel", "Name", "Slot", "Tier", "Dmg")
	return SubHeaderStyle.Render(header[:min(len(header), w.getEffectiveWidth())])
}

// getEffectiveWidth returns the usable width for rendering
func (w *WeaponsListView) getEffectiveWidth() int {
	// Reserve space for borders and padding
	effective := w.Width - 4
	if effective < 80 {
		effective = 80
	}
	return effective
}

// hasMoreItems checks if there are items beyond current page
func (w *WeaponsListView) hasMoreItems() bool {
	return w.CurrentPage < w.TotalPages-1 || (w.SelectedIdx%w.PageSize != 0 && w.SelectedIdx != 0)
}

// renderScrollIndicator shows visual indication of scroll position
func (w *WeaponsListView) renderScrollIndicator() string {
	total := len(w.Weapons)
	_, end := w.getPageRange()
	percent := (end * 100) / total
	filled := "█"
	empty := "░"
	barLen := 20
	filledLen := (percent * barLen) / 100
	bar := strings.Repeat(filled, filledLen) + strings.Repeat(empty, barLen-filledLen)
	return DimStyle.Render(fmt.Sprintf("[%s] %d%%", bar, percent))
}

// renderFilterInfo renders the current search filter status
func (w *WeaponsListView) renderFilterInfo() string {
	if w.isFiltering {
		return HighlightStyle.Render(fmt.Sprintf("Filter: %s_", w.searchInput))
	}
	if w.filterApplied {
		return DimStyle.Render(fmt.Sprintf("Filter: %s", w.searchInput))
	}
	return ""
}

// formatWeaponLine formats a single weapon row with consistent spacing
func (w *WeaponsListView) formatWeaponLine(weapon WeaponData, selected bool) string {
	prefix := " "
	if selected {
		prefix = "▶"
	}

	tierColor := w.getTierColor(weapon.Tier)
	tierText := formatTierCell(weapon.Tier)
	tierStr := tierColor.Render(tierText)

	name := truncateText(weapon.Name, 28)
	//	iconURL := truncateText(weapon.IconUrl, 36)
	slot := truncateText(weapon.Slot, 9)
	//	if iconURL == "" {
	//		iconURL = "(no icon)"
	//}
	if slot == "" {
		slot = "(none)"
	}
	dmgIcon := damageTypeIcon(weapon.DamageType)

	line := fmt.Sprintf("%-3s  %-28s  %-9s  %s %s",
		prefix, name, slot, tierStr, dmgIcon)

	if selected {
		return SelectedStyle.Render(line)
	}
	return NormalStyle.Render(line)
}

func truncateText(value string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(value) <= maxLen {
		return value
	}
	if maxLen <= 3 {
		return value[:maxLen]
	}
	return value[:maxLen-3] + "..."
}

func formatTierCell(tier string) string {
	inner := truncateText(tier, 9)
	label := fmt.Sprintf("[%s]", inner)
	return fmt.Sprintf("%-11s", label)
}

// getTierColor returns the appropriate lipgloss style for a tier
func (w *WeaponsListView) getTierColor(tier string) lipgloss.Style {
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

// getPowerColor returns color based on power level
func (w *WeaponsListView) getPowerColor(power int) lipgloss.Style {
	switch {
	case power >= 800:
		return lipgloss.NewStyle().Foreground(ColorYellow).Bold(true)
	case power >= 750:
		return lipgloss.NewStyle().Foreground(ColorGreen).Bold(true)
	case power >= 700:
		return lipgloss.NewStyle().Foreground(ColorBlue)
	default:
		return lipgloss.NewStyle().Foreground(ColorDim)
	}
}

// Update handles key presses for weapons list
func (w *WeaponsListView) Update(msg tea.Msg) error {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if w.isFiltering {
			w.handleFilterInput(msg)
			return nil
		}

		w.handleNavigationKeys(msg)
	}
	return nil
}

// handleFilterInput processes keyboard input while in filter mode
func (w *WeaponsListView) handleFilterInput(msg tea.KeyMsg) {
	switch msg.String() {
	case "enter":
		w.isFiltering = false
	case "esc":
		w.isFiltering = false
		if !w.filterApplied {
			w.searchInput = ""
		}
	case "backspace":
		if len(w.searchInput) > 0 {
			w.searchInput = w.searchInput[:len(w.searchInput)-1]
			w.applyFilter()
		}
	default:
		if len(msg.Runes) > 0 {
			w.searchInput += string(msg.Runes)
			w.applyFilter()
		}
	}
}

// handleNavigationKeys processes navigation and filter cycling keys
func (w *WeaponsListView) handleNavigationKeys(msg tea.KeyMsg) {
	switch msg.String() {
	case "/":
		w.isFiltering = true
		if !w.filterApplied {
			w.searchInput = ""
		}
	case "r":
		w.cycleRarityFilter()
	case "t":
		w.cycleTypeFilter()
	case "d":
		w.cycleDamageFilter()
	case "up":
		w.moveSelectionUp()
	case "down":
		w.moveSelectionDown()
	case "left":
		w.movePreviousPage()
	case "right":
		w.moveNextPage()
	case "backspace":
		if w.filterApplied {
			w.clearFilter()
		}
	}
}

// moveSelectionUp moves selection up by one item, adjusting page if needed
func (w *WeaponsListView) moveSelectionUp() {
	if w.SelectedIdx > 0 {
		w.SelectedIdx--
		if w.SelectedIdx < w.CurrentPage*w.PageSize {
			w.CurrentPage--
		}
	}
}

// moveSelectionDown moves selection down by one item, adjusting page if needed
func (w *WeaponsListView) moveSelectionDown() {
	if w.SelectedIdx < len(w.Weapons)-1 {
		w.SelectedIdx++
		if w.SelectedIdx >= (w.CurrentPage+1)*w.PageSize {
			w.CurrentPage++
		}
	}
}

// movePreviousPage navigates to previous page
func (w *WeaponsListView) movePreviousPage() {
	if w.CurrentPage > 0 {
		w.CurrentPage--
		w.SelectedIdx = w.CurrentPage * w.PageSize
	}
}

// moveNextPage navigates to next page
func (w *WeaponsListView) moveNextPage() {
	if w.CurrentPage < w.TotalPages-1 {
		w.CurrentPage++
		w.SelectedIdx = w.CurrentPage * w.PageSize
		if w.SelectedIdx >= len(w.Weapons) {
			w.SelectedIdx = len(w.Weapons) - 1
		}
	}
}

// IsFiltering returns whether the user is currently entering a search query
func (w *WeaponsListView) IsFiltering() bool {
	return w.isFiltering
}

// applyFilter applies all active filters to the weapons list
func (w *WeaponsListView) applyFilter() {
	query := strings.TrimSpace(strings.ToLower(w.searchInput))
	rarityFilter := strings.ToLower(w.currentRarityFilter())
	typeFilter := strings.ToLower(w.currentTypeFilter())
	damageFilter := strings.ToLower(w.currentDamageFilter())

	filtered := make([]WeaponData, 0, len(w.allWeapons))
	for _, weapon := range w.allWeapons {
		if w.matchesFilters(weapon, query, rarityFilter, typeFilter, damageFilter) {
			filtered = append(filtered, weapon)
		}
	}

	w.Weapons = filtered
	w.filterApplied = w.hasActiveFilter()
	w.SelectedIdx = 0
	w.CurrentPage = 0
	w.calculatePages()
}

// matchesFilters checks if a weapon matches all currently active filters
func (w *WeaponsListView) matchesFilters(weapon WeaponData, query, rarityFilter, typeFilter, damageFilter string) bool {
	if rarityFilter != "all" && strings.ToLower(strings.TrimSpace(weapon.Tier)) != rarityFilter {
		return false
	}
	if typeFilter != "all" && strings.ToLower(strings.TrimSpace(weapon.Type)) != typeFilter {
		return false
	}
	if damageFilter != "all" && strings.ToLower(normalizedDamageType(weapon.DamageType)) != damageFilter {
		return false
	}
	if query != "" {
		haystack := strings.ToLower(weapon.Name + " " + weapon.Type + " " + weapon.Tier + " " + weapon.Slot)
		return strings.Contains(haystack, query)
	}
	return true
}

func (w *WeaponsListView) clearFilter() {
	w.Weapons = append([]WeaponData(nil), w.allWeapons...)
	w.filterApplied = false
	w.SelectedIdx = 0
	w.CurrentPage = 0
	w.calculatePages()
}

// refreshTypeOptions rebuilds the list of available weapon types from allWeapons
func (w *WeaponsListView) refreshTypeOptions() {
	typeSet := map[string]struct{}{}
	for _, weapon := range w.allWeapons {
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

	w.typeOptions = append([]string{"All"}, types...)
	if w.filterState.TypeIdx >= len(w.typeOptions) {
		w.filterState.TypeIdx = 0
	}
}

// getFilterValue safely retrieves a filter value by index from options
func (w *WeaponsListView) getFilterValue(idx int, options []string) string {
	if idx < 0 || idx >= len(options) {
		return "All"
	}
	return options[idx]
}

// currentRarityFilter returns the currently selected rarity filter
func (w *WeaponsListView) currentRarityFilter() string {
	return w.getFilterValue(w.filterState.RarityIdx, w.rarityOptions)
}

// currentTypeFilter returns the currently selected type filter
func (w *WeaponsListView) currentTypeFilter() string {
	return w.getFilterValue(w.filterState.TypeIdx, w.typeOptions)
}

// currentDamageFilter returns the currently selected damage filter
func (w *WeaponsListView) currentDamageFilter() string {
	return w.getFilterValue(w.filterState.DamageIdx, w.damageOptions)
}

// cycleRarityFilter moves to the next rarity filter option
func (w *WeaponsListView) cycleRarityFilter() {
	w.filterState.RarityIdx = (w.filterState.RarityIdx + 1) % len(w.rarityOptions)
	w.applyFilter()
}

// cycleTypeFilter moves to the next type filter option
func (w *WeaponsListView) cycleTypeFilter() {
	w.filterState.TypeIdx = (w.filterState.TypeIdx + 1) % len(w.typeOptions)
	w.applyFilter()
}

// cycleDamageFilter moves to the next damage filter option
func (w *WeaponsListView) cycleDamageFilter() {
	w.filterState.DamageIdx = (w.filterState.DamageIdx + 1) % len(w.damageOptions)
	w.applyFilter()
}

// resetFilters clears all filter selections back to "All"
func (w *WeaponsListView) resetFilters() {
	w.filterState = FilterState{RarityIdx: 0, TypeIdx: 0, DamageIdx: 0}
	w.searchInput = ""
	w.applyFilter()
}

// hasActiveFilter checks if any filter (search, rarity, type, or damage) is active
func (w *WeaponsListView) hasActiveFilter() bool {
	return strings.TrimSpace(w.searchInput) != "" ||
		!strings.EqualFold(w.currentRarityFilter(), "All") ||
		!strings.EqualFold(w.currentTypeFilter(), "All") ||
		!strings.EqualFold(w.currentDamageFilter(), "All")
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Init initializes the view
func (w *WeaponsListView) Init() tea.Cmd {
	return nil
}
