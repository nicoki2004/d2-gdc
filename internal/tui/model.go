package tui

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nicoki2004/g2-drc/internal/repository"
)

// Model is the main bubbletea model for the TUI
type Model struct {
	// Global state
	repo        repository.WeaponRepository
	currentView ViewType

	// Cached data
	allWeapons  []WeaponData
	characters  []CharacterData
	selectedIdx int
	sortBy      SortBy

	// View models
	homeView   *HomeView
	listView   *WeaponsListView
	detailView *WeaponDetailView
	dupView    *DuplicateView
	compView   *ComparisonView
	searchView *GodRollView

	// Search state
	searchQuery string
	filters     SearchFilters

	// UI state
	width, height int
	err           error
	loading       bool
	message       string
	statsCache    map[string]map[string]int
	perksCache    map[string][]PerkData
}

// NewModel creates a new TUI model
func NewModel(repo repository.WeaponRepository) Model {
	return Model{
		repo:        repo,
		currentView: ViewHome,
		filters: SearchFilters{
			MinPower: 0,
			MaxPower: 999,
		},
		sortBy:     SortByName,
		loading:    true,
		statsCache: make(map[string]map[string]int),
		perksCache: make(map[string][]PerkData),
	}
}

// Init initializes the model and loads initial data
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadWeapons(),
		m.loadCharacters(),
	)
}

// Update handles key presses and messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case weaponsLoadedMsg:
		m.allWeapons = msg.weapons
		m.loading = false
		m.message = fmt.Sprintf("Loaded %d weapons", len(msg.weapons))
		m.statsCache = make(map[string]map[string]int)
		m.perksCache = make(map[string][]PerkData)
		if m.listView == nil {
			m.listView = NewWeaponsListView(m.allWeapons, m.width-4, m.height-10)
		} else {
			m.listView.SetWeapons(m.allWeapons)
		}
		if m.dupView != nil {
			m.dupView.SetGroups(m.groupWeaponsByHash())
		}
		if m.compView != nil {
			m.compView.SetWeapons(m.allWeapons)
		}
		return m, nil
	case charactersLoadedMsg:
		m.characters = msg.characters
		m.message = fmt.Sprintf("Loaded %d characters", len(msg.characters))
		return m, nil
	case errorMsg:
		m.err = msg.err
		m.message = fmt.Sprintf("Error: %v", msg.err)
		return m, nil
	}
	return m, nil
}

// View renders the current view
func (m Model) View() string {
	if m.loading {
		return HeaderStyle.Render("Loading arsenal...")
	}

	if m.err != nil {
		return ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}

	switch m.currentView {
	case ViewHome:
		return m.renderHome()
	case ViewList:
		return m.renderList()
	case ViewDetail:
		return m.renderDetail()
	case ViewDuplicates:
		return m.renderDuplicates()
	case ViewComparison:
		return m.renderComparison()
	case ViewGodRoll:
		return m.renderGodRoll()
	case ViewHelp:
		return m.renderHelp()
	default:
		return m.renderHome()
	}
}

// Handle key presses
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global hotkeys
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		if m.currentView == ViewList && m.listView != nil && m.listView.IsFiltering() {
			m.listView.Update(msg)
			return m, nil
		}
		if m.currentView == ViewComparison && m.compView != nil && m.compView.IsFiltering() {
			m.compView.Update(msg)
			return m, nil
		}
		if m.currentView != ViewHome {
			m.currentView = ViewHome
		}
		return m, nil
	}

	// Main menu navigation
	if m.currentView == ViewHome {
		switch msg.String() {
		case "up":
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		case "down":
			if m.selectedIdx < 5 {
				m.selectedIdx++
			}
		case "enter":
			return m.openHomeSelection()
		case "1":
			m.selectedIdx = 0
			return m.openHomeSelection()
		case "2":
			m.selectedIdx = 1
			return m.openHomeSelection()
		case "3":
			m.selectedIdx = 2
			return m.openHomeSelection()
		case "4":
			m.selectedIdx = 3
			return m.openHomeSelection()
		case "h":
			m.selectedIdx = 4
			return m.openHomeSelection()
		case "home":
			m.currentView = ViewHome
		}
		return m, nil
	}

	// Pass key presses to current view
	switch m.currentView {
	case ViewList:
		if m.listView != nil {
			m.listView.Update(msg)
			// Check for Enter key on weapon selection
			if msg.String() == "enter" && !m.listView.IsFiltering() && m.listView.SelectedIdx < len(m.listView.Weapons) {
				m.currentView = ViewDetail
				weapon := m.listView.Weapons[m.listView.SelectedIdx]
				weapon = m.enrichWeaponForDetail(weapon)
				if m.detailView == nil {
					m.detailView = NewWeaponDetailView(weapon)
				} else {
					m.detailView.SetWeapon(weapon)
				}
			}
		}
	case ViewDetail:
		if m.detailView != nil {
			m.detailView.Update(msg)
		}
		// Press Enter to go back to list (ESC remains global to home)
		if msg.String() == "enter" || msg.String() == "backspace" {
			m.currentView = ViewList
		}
	case ViewDuplicates:
		if m.dupView != nil {
			m.dupView.Update(msg)
		}
	case ViewComparison:
		if m.compView != nil {
			m.compView.Update(msg)
			if msg.String() == "enter" || msg.String() == " " {
				if current, ok := m.compView.CurrentWeapon(); ok && m.compView.IsSelected(current.InstanceID) {
					m.compView.UpsertSelectedWeapon(m.enrichWeaponForComparison(current))
				}
			}
		}
	case ViewGodRoll:
		if m.searchView != nil {
			m.searchView.Update(msg)
		}
	}

	return m, nil
}

func (m Model) openHomeSelection() (tea.Model, tea.Cmd) {
	switch m.selectedIdx {
	case 0:
		m.currentView = ViewList
		if m.listView == nil {
			m.listView = NewWeaponsListView(m.allWeapons, m.width-4, m.height-10)
		}
	case 1:
		m.currentView = ViewDuplicates
		if m.dupView == nil {
			m.dupView = NewDuplicateView(m.groupWeaponsByHash())
		} else {
			m.dupView.SetGroups(m.groupWeaponsByHash())
		}
	case 2:
		m.currentView = ViewComparison
		if m.compView == nil {
			m.compView = NewComparisonView(4)
		}
		m.compView.SetWeapons(m.allWeapons)
	case 3:
		m.currentView = ViewGodRoll
		if m.searchView == nil {
			m.searchView = NewGodRollView()
		}
	case 4:
		m.currentView = ViewHelp
	case 5:
		return m, tea.Quit
	}
	return m, nil
}

// groupWeaponsByHash groups weapons by hash for the duplicate viewer
func (m Model) groupWeaponsByHash() [][]WeaponData {
	hashMap := make(map[int64][]WeaponData)

	for _, weapon := range m.allWeapons {
		weapon.Stats = m.getWeaponStatsForDuplicate(weapon.InstanceID)
		weapon.Location = m.formatWeaponLocation(weapon)
		hashMap[weapon.Hash] = append(hashMap[weapon.Hash], weapon)
	}

	var groups [][]WeaponData
	for _, group := range hashMap {
		if len(group) > 1 {
			sort.Slice(group, func(i, j int) bool {
				if group[i].Power == group[j].Power {
					return group[i].Name < group[j].Name
				}
				return group[i].Power > group[j].Power
			})
			groups = append(groups, group)
		}
	}

	sort.Slice(groups, func(i, j int) bool {
		if len(groups[i]) == len(groups[j]) {
			iTopPower, jTopPower := 0, 0
			if len(groups[i]) > 0 {
				iTopPower = groups[i][0].Power
			}
			if len(groups[j]) > 0 {
				jTopPower = groups[j][0].Power
			}
			if iTopPower == jTopPower {
				return groups[i][0].Name < groups[j][0].Name
			}
			return iTopPower > jTopPower
		}
		return len(groups[i]) > len(groups[j])
	})

	return groups
}

func (m Model) getWeaponStatsForDuplicate(instanceID string) map[string]int {
	if stats, ok := m.statsCache[instanceID]; ok {
		return stats
	}

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	stats, err := m.repo.GetWeaponStats(ctx, instanceID)
	if err != nil {
		stats = map[string]int{}
	}
	m.statsCache[instanceID] = stats
	return stats
}

func (m Model) getWeaponPerksForComparison(instanceID string) []PerkData {
	if perks, ok := m.perksCache[instanceID]; ok {
		return perks
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	rows, err := m.repo.GetWeaponPerks(ctx, instanceID)
	if err != nil {
		m.perksCache[instanceID] = []PerkData{}
		return []PerkData{}
	}

	perks := make([]PerkData, 0, len(rows))
	for _, p := range rows {
		perks = append(perks, PerkData{
			Hash:       uint32(p.Hash),
			Name:       p.Name,
			IsEquipped: p.IsEquipped,
			SocketIdx:  p.SocketIdx,
		})
	}

	m.perksCache[instanceID] = perks
	return perks
}

// Render functions for each view
func (m Model) renderHome() string {
	title := HeaderStyle.Render("🔫 DESTINY 2 GOD ROLL CHECKER")
	menuItems := []string{
		"[1] Search Weapons",
		"[2] Find Duplicates",
		"[3] Compare Weapons",
		"[4] Search God Rolls",
		"[h] Help",
		"[q] Quit",
	}

	lines := []string{SubHeaderStyle.Render("Main Menu"), ""}
	for i, item := range menuItems {
		marker := " "
		style := NormalStyle
		if i == m.selectedIdx {
			marker = "▶"
			style = SelectedStyle
		}
		lines = append(lines, style.Render(fmt.Sprintf("%s %s", marker, item)))
	}

	footer := fmt.Sprintf("\n%s\nTotal Weapons: %d | Characters: %d",
		DimStyle.Render("Arrow keys to move, Enter to select"),
		len(m.allWeapons),
		len(m.characters),
	)

	return JoinVertical(title, BorderBoxStyle.Render(strings.Join(lines, "\n")), footer)
}

func (m Model) renderList() string {
	if m.listView == nil {
		return "Loading weapon list..."
	}
	return m.listView.View()
}

func (m Model) renderDetail() string {
	if m.detailView != nil {
		return m.detailView.View()
	}
	if m.listView == nil || m.listView.SelectedIdx >= len(m.listView.Weapons) {
		return "No weapon selected"
	}
	weapon := m.listView.Weapons[m.listView.SelectedIdx]
	return m.renderWeaponDetail(weapon)
}

func (m Model) enrichWeaponForDetail(weapon WeaponData) WeaponData {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stats, err := m.repo.GetWeaponStats(ctx, weapon.InstanceID)
	if err == nil {
		weapon.Stats = stats
	}

	perks, err := m.repo.GetWeaponPerks(ctx, weapon.InstanceID)
	if err == nil {
		weapon.Perks = make([]PerkData, 0, len(perks))
		for _, p := range perks {
			weapon.Perks = append(weapon.Perks, PerkData{
				Hash:       uint32(p.Hash),
				Name:       p.Name,
				IsEquipped: p.IsEquipped,
				SocketIdx:  p.SocketIdx,
			})
		}
	}

	weapon.Location = m.formatWeaponLocation(weapon)

	return weapon
}

func (m Model) enrichWeaponForComparison(weapon WeaponData) WeaponData {
	weapon.Stats = m.getWeaponStatsForDuplicate(weapon.InstanceID)
	weapon.Perks = m.getWeaponPerksForComparison(weapon.InstanceID)
	weapon.Location = m.formatWeaponLocation(weapon)
	return weapon
}

func (m Model) formatWeaponLocation(weapon WeaponData) string {
	raw := strings.TrimSpace(weapon.Location)
	if raw == "" {
		return "(unknown)"
	}

	if strings.EqualFold(raw, "Vault") {
		return "Vault"
	}

	prefix, charID, ok := strings.Cut(raw, ":")
	if !ok {
		return raw
	}

	normalizedPrefix := strings.ToLower(strings.TrimSpace(prefix))
	charID = strings.TrimSpace(charID)
	if charID == "" {
		return raw
	}

	charLabel := charID
	for _, c := range m.characters {
		if c.CharacterID == charID {
			charLabel = fmt.Sprintf("%s (%s)", classTypeName(c.ClassType), shortID(c.CharacterID))
			break
		}
	}

	switch normalizedPrefix {
	case "equipped":
		return charLabel
	case "inventory":
		return "Inventory: " + charLabel
	default:
		return raw
	}
}

func classTypeName(classType int) string {
	switch classType {
	case 0:
		return "Titan"
	case 1:
		return "Hunter"
	case 2:
		return "Warlock"
	default:
		return "Unknown"
	}
}

func shortID(id string) string {
	if len(id) <= 8 {
		return id
	}
	return id[:8] + "..."
}

func (m Model) renderWeaponDetail(w WeaponData) string {
	title := HeaderStyle.Render(fmt.Sprintf("⚔️ %s [%s]", w.Name, w.Tier))

	stats := fmt.Sprintf(`
Power: %d | Kills: %d | Level: %d
Type: %s | Slot: %s | Damage: %s
Location: %s
		`, w.Power, w.Kills, w.Level, w.Type, w.Slot, w.DamageType, w.Location)

	return JoinVertical(title, stats)
}

func (m Model) renderDuplicates() string {
	title := HeaderStyle.Render("🔄 DUPLICATE WEAPONS")
	if m.dupView == nil {
		return JoinVertical(title, "Loading duplicates...")
	}
	return m.dupView.View()
}

func (m Model) renderComparison() string {
	title := HeaderStyle.Render("⚖️ WEAPON COMPARISON")
	if m.compView == nil {
		return JoinVertical(title, "No weapons selected for comparison")
	}
	return m.compView.View()
}

func (m Model) renderGodRoll() string {
	title := HeaderStyle.Render("✨ GOD ROLL SEARCH")
	if m.searchView == nil {
		return JoinVertical(title, "Loading god roll data...")
	}
	return m.searchView.View()
}

func (m Model) renderHelp() string {
	help := `
🔫 DESTINY 2 GOD ROLL CHECKER - Help

KEYBINDS:
  1           - Search Weapons
  2           - Find Duplicates
  3           - Compare Weapons
  4           - Search God Rolls
  h           - Show this help
  q/ctrl+c    - Quit
  esc         - Go back to home
  ↑↓          - Navigate items
  →←          - Change page
  Enter       - Select/View details

FEATURES:
  • Search weapons by name or attributes
  • View duplicate weapons grouped by hash
  • Compare up to 4 weapons side-by-side
  • Find god rolls matching your criteria
  • Filter by power, rarity, slot, etc.

Press any key to return to menu...
	`
	return BorderStyle.Render(help)
}

// Command functions
func (m Model) loadWeapons() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		weapons, err := m.repo.GetAllWeapons(ctx)
		if err != nil {
			return errorMsg{err}
		}

		return weaponsLoadedMsg{weapons: convertToWeaponData(weapons)}
	}
}

func (m Model) loadCharacters() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		characters, err := m.repo.GetAllCharacters(ctx)
		if err != nil {
			return errorMsg{err}
		}

		return charactersLoadedMsg{characters: convertToCharacterData(characters)}
	}
}

// Message types
type weaponsLoadedMsg struct {
	weapons []WeaponData
}

type charactersLoadedMsg struct {
	characters []CharacterData
}

type errorMsg struct {
	err error
}

type weaponSelectedMsg struct {
	weapon WeaponData
}
