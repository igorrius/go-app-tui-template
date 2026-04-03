package screen

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/igorrius/go-app-tui-template/internal/logging"
	"github.com/igorrius/go-app-tui-template/internal/tui/nav"
)

// New creates a screen model for the given identifier.
// logCh is required for LogViewID.
func New(id nav.ID, logCh <-chan logging.LogEvent) (tea.Model, error) {
	switch id {
	case nav.DashboardID:
		return NewDashboard(), nil
	case nav.LogViewID:
		return NewLogView(logCh), nil
	default:
		return nil, fmt.Errorf("unknown screen id: %s", id)
	}
}
