package screen

import (
	"log/slog"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/igorrius/go-app-tui-template/internal/logging"
)

// LogView displays live log entries as plain text lines inside the TUI.
// It receives LogEvent messages from a provided channel and appends
// each formatted entry to an in-memory buffer, showing only the last height lines.
type LogView struct {
	formatter logging.ColorTextFormatter
	lines     []string
	ch        <-chan logging.LogEvent
	width     int
	height    int
}

// NewLogView creates a new LogView screen with the given log event channel.
func NewLogView(ch <-chan logging.LogEvent) *LogView {
	return &LogView{ch: ch}
}

// Init emits a subscription info log and returns the Cmd to wait for the next event.
func (lv *LogView) Init() tea.Cmd {
	slog.Debug("LogView initialized, subscribed to log events", "module", "TUI_LOGVIEW")
	return lv.waitForEvent()
}

// waitForEvent returns a Cmd that blocks until the next LogEvent arrives on the channel.
func (lv *LogView) waitForEvent() tea.Cmd {
	ch := lv.ch
	return func() tea.Msg {
		return <-ch
	}
}

// Update handles window resize and incoming LogEvent messages.
func (lv *LogView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		lv.width = msg.Width
		lv.height = msg.Height
		lv.trimBuffer()

	case logging.LogEvent:
		line := lv.formatter.Format(msg)
		lv.lines = append(lv.lines, line)
		lv.trimBuffer()
		return lv, lv.waitForEvent()
	}

	return lv, nil
}

// View renders the buffered log lines as plain text.
func (lv *LogView) View() tea.View {
	return tea.NewView(strings.Join(lv.lines, "\n"))
}

// trimBuffer trims the line buffer to at most height entries.
func (lv *LogView) trimBuffer() {
	if lv.height > 0 && len(lv.lines) > lv.height {
		lv.lines = lv.lines[len(lv.lines)-lv.height:]
	}
}
