package screen

import (
	"log/slog"
	"strings"

	tea "charm.land/bubbletea/v2"
)

// Dashboard displays an "In construction" placeholder centered in the available area.
type Dashboard struct {
	width  int
	height int
}

// NewDashboard creates a new Dashboard screen.
func NewDashboard() *Dashboard {
	return &Dashboard{}
}

func (d *Dashboard) Init() tea.Cmd {
	slog.Debug("Dashboard initialized", "module", "TUI_DASHBOARD")
	return nil
}

func (d *Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		d.width = msg.Width
		d.height = msg.Height
	}

	return d, nil
}

func (d *Dashboard) View() tea.View {
	if d.width == 0 || d.height == 0 {
		return tea.NewView("")
	}

	text := "In construction"

	// Center horizontally
	pad := (d.width - len(text)) / 2
	if pad < 0 {
		pad = 0
	}
	line := strings.Repeat(" ", pad) + text

	// Center vertically while accounting for the single content line.
	topPad := (d.height - 1) / 2
	if topPad < 0 {
		topPad = 0
	}
	var b strings.Builder
	for i := 0; i < topPad; i++ {
		b.WriteString("\n")
	}
	b.WriteString(line)

	return tea.NewView(b.String())
}
