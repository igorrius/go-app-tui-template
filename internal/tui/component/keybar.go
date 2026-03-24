package component

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/igorrius/go-app-tui-template/internal/tui/screen"
)

type keyBarItem struct {
	key      string
	label    string
	screenID screen.ID
}

var keyBarItems = []keyBarItem{
	{key: "F1", label: "Dashboard", screenID: screen.DashboardID},
	{key: "F9", label: "Logs", screenID: screen.LogViewID},
	{key: "F10", label: "Quit"},
}

var (
	keyStyle         = lipgloss.NewStyle().Bold(true)
	labelStyle       = lipgloss.NewStyle()
	activeLabelStyle = lipgloss.NewStyle().Reverse(true)
)

// KeyBar renders the function key bar at the bottom of the TUI.
type KeyBar struct {
	width    int
	items    []keyBarItem
	activeID screen.ID
}

// NewKeyBar creates a key bar component.
func NewKeyBar() *KeyBar {
	return &KeyBar{
		items:    keyBarItems,
		activeID: screen.DashboardID,
	}
}

func (k *KeyBar) Init() tea.Cmd {
	return nil
}

func (k *KeyBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		k.width = msg.Width
	}

	return k, nil
}

func (k *KeyBar) View() tea.View {
	var parts []string
	for _, item := range k.items {
		var rendered string
		if item.screenID != "" && item.screenID == k.activeID {
			rendered = keyStyle.Render(item.key) + " " + activeLabelStyle.Render(item.label)
		} else {
			rendered = keyStyle.Render(item.key) + " " + labelStyle.Render(item.label)
		}
		parts = append(parts, rendered)
	}

	content := strings.Join(parts, "  ")
	v := tea.NewView(lipgloss.NewStyle().Width(k.width).Render(content))
	return v
}

func (k *KeyBar) SetSize(width, _ int) {
	k.width = width
}

func (k *KeyBar) Height() int {
	return 1
}

func (k *KeyBar) SetActiveID(id screen.ID) {
	k.activeID = id
}
