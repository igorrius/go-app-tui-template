package screen

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/igorrius/go-app-tui-template/internal/logging"
)

// ID identifies a routed screen in the TUI.
type ID string

const (
	DashboardID ID = "dashboard"
	LogViewID   ID = "logview"
)

// New creates a screen model for the given identifier.
// logCh is required for LogViewID.
func New(id ID, logCh <-chan logging.LogEvent) (tea.Model, error) {
	switch id {
	case DashboardID:
		return NewDashboard(), nil
	case LogViewID:
		return NewLogView(logCh), nil
	default:
		return nil, fmt.Errorf("unknown screen id: %s", id)
	}
}
