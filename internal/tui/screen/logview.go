package screen

import (
	"log/slog"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/igorrius/go-app-tui-template/internal/logging"
	"github.com/igorrius/go-app-tui-template/internal/tui/component"
)

// LogView displays live log entries as plain text lines inside the TUI.
// It receives LogEvent messages from a provided channel and appends
// each formatted entry to an in-memory buffer, showing only the last height lines.
type LogView struct {
	formatter    logging.ColorTextFormatter
	lines        []string
	masterEvents []logging.LogEvent
	filterBar    *component.FilterBar
	ch           <-chan logging.LogEvent
	width        int
	height       int
}

// NewLogView creates a new LogView screen with the given log event channel.
func NewLogView(ch <-chan logging.LogEvent) *LogView {
	return &LogView{
		ch:        ch,
		filterBar: component.NewFilterBar(),
	}
}

// Init emits a subscription info log.
func (lv *LogView) Init() tea.Cmd {
	slog.Debug("LogView initialized, subscribed to log events", "module", "TUI_LOGVIEW")
	return nil
}

// Update handles window resize and incoming LogEvent messages.
func (lv *LogView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		lv.width = msg.Width
		lv.height = msg.Height
		lv.rebuildBuffer()

	case logging.LogEvent:
		lv.masterEvents = append(lv.masterEvents, msg)
		if len(lv.masterEvents) > 1000 {
			lv.masterEvents = lv.masterEvents[len(lv.masterEvents)-1000:]
		}

		if lv.passesFilter(msg) {
			line := lv.formatter.Format(msg)
			lv.lines = append(lv.lines, line)
			lv.trimBuffer()
		}

	case tea.KeyMsg:
		if lv.filterBar.IsSearching {
			lv.filterBar.HandleKey(msg)
			lv.rebuildBuffer()
			return lv, nil
		}

		k := msg.String()
		switch k {
		case "1", "2", "3", "4":
			levels := []slog.Level{slog.LevelDebug, slog.LevelInfo, slog.LevelWarn, slog.LevelError}
			idx := int(k[0] - '1')
			lv.filterBar.ToggleLevel(levels[idx])
			lv.rebuildBuffer()
		case "/":
			lv.filterBar.IsSearching = true
			return lv, nil
		case "esc":
			if lv.filterBar.SearchQuery != "" {
				lv.filterBar.SearchQuery = ""
				lv.rebuildBuffer()
			}
		}
	}

	return lv, nil
}

// View renders the buffered log lines as plain text.
func (lv *LogView) View() tea.View {
	contentHeight := lv.contentLinesHeight()
	var b strings.Builder

	for i := 0; i < contentHeight; i++ {
		if i < len(lv.lines) {
			b.WriteString(lv.lines[i])
		}
		if i < contentHeight-1 {
			b.WriteByte('\n')
		}
	}

	filterView := lv.filterBar.View()
	content := b.String()
	// Always ensure there's a newline before the filter bar if logs exist.
	if content != "" {
		content += "\n"
	}
	content += filterView.Content
	// Do NOT add a trailing newline here; App.View will handle it.

	return tea.NewView(content)
}

func (lv *LogView) rebuildBuffer() {
	lv.lines = nil
	for _, e := range lv.masterEvents {
		if lv.passesFilter(e) {
			lv.lines = append(lv.lines, lv.formatter.Format(e))
		}
	}
	lv.trimBuffer()
}

func (lv *LogView) passesFilter(e logging.LogEvent) bool {
	// Level filter
	if !lv.filterBar.EnabledLevels[e.Level] {
		return false
	}

	// Text filter
	if lv.filterBar.SearchQuery == "" {
		return true
	}

	q := strings.ToLower(lv.filterBar.SearchQuery)
	if strings.Contains(strings.ToLower(e.Message), q) {
		return true
	}

	for _, a := range e.Attrs {
		if strings.Contains(strings.ToLower(a.Key), q) || strings.Contains(strings.ToLower(a.Value.String()), q) {
			return true
		}
	}

	return false
}

// trimBuffer trims the line buffer to at most height-1 entries (reserved for filter bar).
func (lv *LogView) trimBuffer() {
	h := lv.contentLinesHeight()
	if h > 0 && len(lv.lines) > h {
		lv.lines = lv.lines[len(lv.lines)-h:]
	}
}

func (lv *LogView) contentLinesHeight() int {
	if lv.height <= 1 {
		return 0
	}
	// Filter bar takes 1 line.
	return lv.height - 1
}
