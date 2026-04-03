package component

import (
	"fmt"
	"log/slog"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

var (
	checkedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
	uncheckedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))  // Red
	filterLabelStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	searchInputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("229"))
)

type FilterBar struct {
	width         int
	EnabledLevels map[slog.Level]bool
	SearchQuery   string
	IsSearching   bool
}

func NewFilterBar() *FilterBar {
	return &FilterBar{
		EnabledLevels: map[slog.Level]bool{
			slog.LevelDebug: true,
			slog.LevelInfo:  true,
			slog.LevelWarn:  true,
			slog.LevelError: true,
		},
	}
}

func (f *FilterBar) Init() tea.Cmd {
	return nil
}

func (f *FilterBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		f.width = msg.Width
	}
	return f, nil
}

func (f *FilterBar) HandleKey(msg tea.KeyMsg) {
	if !f.IsSearching {
		return
	}

	switch msg.String() {
	case "enter":
		f.IsSearching = false
	case "esc":
		f.SearchQuery = ""
		f.IsSearching = false
	case "backspace":
		if len(f.SearchQuery) > 0 {
			f.SearchQuery = f.SearchQuery[:len(f.SearchQuery)-1]
		}
	default:
		if len(msg.String()) == 1 {
			f.SearchQuery += msg.String()
		}
	}
}

func (f *FilterBar) View() tea.View {
	var b strings.Builder

	// Levels
	levels := []struct {
		key   int
		level slog.Level
		name  string
	}{
		{1, slog.LevelDebug, "DEBUG"},
		{2, slog.LevelInfo, "INFO"},
		{3, slog.LevelWarn, "WARN"},
		{4, slog.LevelError, "ERROR"},
	}

	for i, l := range levels {
		var check string
		if f.EnabledLevels[l.level] {
			check = checkedStyle.Render("[x]")
		} else {
			check = uncheckedStyle.Render("[ ]")
		}
		fmt.Fprintf(&b, "%d%s %-5s", l.key, check, l.name)
		if i < len(levels)-1 {
			b.WriteString("  ")
		}
	}

	// Separator
	if !f.IsSearching && f.SearchQuery == "" {
		b.WriteString("  " + lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("Press / to filter text"))
	} else {
		b.WriteString("  |  ")
	}

	// Search
	if f.IsSearching || f.SearchQuery != "" {
		label := "Filter: "
		b.WriteString(filterLabelStyle.Render(label))

		query := f.SearchQuery
		if f.IsSearching {
			if query == "" {
				b.WriteString(searchInputStyle.Render("<typing...>"))
			} else {
				b.WriteString(searchInputStyle.Render(query))
			}
		} else {
			b.WriteString(searchInputStyle.Render(query))
		}
	}

	content := b.String()
	style := lipgloss.NewStyle().Width(f.width)
	if f.IsSearching {
		style = style.Background(lipgloss.Color("235")) // Subtle dark background
	}

	return tea.NewView(style.Render(content))
}

func (f *FilterBar) ToggleLevel(l slog.Level) {
	f.EnabledLevels[l] = !f.EnabledLevels[l]
}

func (f *FilterBar) SetSearch(query string) {
	f.SearchQuery = query
}
