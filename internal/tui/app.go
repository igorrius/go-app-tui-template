package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/igorrius/go-app-tui-template/internal/config"
	"github.com/igorrius/go-app-tui-template/internal/logging"
	"github.com/igorrius/go-app-tui-template/internal/tui/component"
	"github.com/igorrius/go-app-tui-template/internal/tui/nav"
	"github.com/igorrius/go-app-tui-template/internal/tui/screen"
	vein "github.com/igorrius/go-vein"
)

// App is the root bubbletea model managing active screen, terminal size, and key bar.
type App struct {
	cfg            *config.AppConfig
	dispatcher     *vein.Dispatcher
	logCh          <-chan logging.LogEvent
	screens        map[nav.ID]tea.Model
	activeScreenID nav.ID
	screen         tea.Model
	keyBar         *component.KeyBar
	width          int
	height         int
}

// NewApp creates a new root TUI application model.
func NewApp(cfg *config.AppConfig, dispatcher *vein.Dispatcher) *App {
	sub := vein.SubscribeTo[logging.LogEvent](dispatcher)
	logCh := sub.OnC()

	screens := make(map[nav.ID]tea.Model)
	scrDashboard, _ := screen.New(nav.DashboardID, logCh)
	screens[nav.DashboardID] = scrDashboard

	scrLogView, _ := screen.New(nav.LogViewID, logCh)
	screens[nav.LogViewID] = scrLogView

	return &App{
		cfg:            cfg,
		dispatcher:     dispatcher,
		logCh:          logCh,
		screens:        screens,
		activeScreenID: nav.DashboardID,
		screen:         screens[nav.DashboardID],
		keyBar:         component.NewKeyBar(),
	}
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.screen.Init(),
		a.keyBar.Init(),
		a.waitForLog(),
	)
}

func (a *App) waitForLog() tea.Cmd {
	return func() tea.Msg {
		return <-a.logCh
	}
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

		sizeMsg := tea.WindowSizeMsg{
			Width:  a.width,
			Height: a.contentHeight(),
		}
		for id, scr := range a.screens {
			nextScr, _ := scr.Update(sizeMsg)
			a.screens[id] = nextScr
		}
		a.screen = a.screens[a.activeScreenID]
		a.keyBar.SetSize(a.width, a.keyBar.Height())
		return a, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "f10", "ctrl+c":
			return a, tea.Quit
		case "f1":
			return a, a.routeTo(nav.DashboardID)
		case "f9":
			return a, a.routeTo(nav.LogViewID)
		}

	case logging.LogEvent:
		cmds = append(cmds, a.waitForLog())
	}

	// Broadcast global messages to all screens to persist state
	switch msg.(type) {
	case logging.LogEvent:
		for id, scr := range a.screens {
			nextScr, _ := scr.Update(msg)
			a.screens[id] = nextScr
		}
		a.screen = a.screens[a.activeScreenID]
	default:
		nextScreen, screenCmd := a.screen.Update(msg)
		a.screen = nextScreen
		a.screens[a.activeScreenID] = nextScreen
		if screenCmd != nil {
			cmds = append(cmds, screenCmd)
		}
	}

	nextKeyBar, keyBarCmd := a.keyBar.Update(msg)
	if kb, ok := nextKeyBar.(*component.KeyBar); ok {
		a.keyBar = kb
	}
	if keyBarCmd != nil {
		cmds = append(cmds, keyBarCmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) View() tea.View {
	v := tea.NewView("App")
	v.AltScreen = true

	if a.width == 0 || a.height == 0 {
		return v
	}

	screenView := a.screen.View()
	keyBarView := a.keyBar.View()

	// Screen content lines
	screenContent := screenView.Content
	screenLines := strings.Split(screenContent, "\n")
	screenLineCount := len(screenLines)
	if screenContent == "" {
		screenLineCount = 0
	}

	contentHeight := a.contentHeight()

	var b strings.Builder
	b.WriteString(screenContent)

	// Determine if we need to add a newline and how much padding.
	// We want the KeyBar to be exactly at the last line (index a.height-1).
	// Total lines (including screen + padding + keyBar) should be exactly a.height.

	// If the screen content didn't end with a newline, and we have more lines to add (padding or keybar),
	// we must add a newline to move to the next logical line.
	if screenContent != "" && !strings.HasSuffix(screenContent, "\n") {
		b.WriteByte('\n')
	}

	// Remaining lines to fill before the KeyBar.
	// contentHeight is height - 1.
	padLines := contentHeight - screenLineCount
	for i := 0; i < padLines; i++ {
		b.WriteByte('\n')
	}

	b.WriteString(keyBarView.Content)

	v.SetContent(b.String())
	return v
}

func (a *App) routeTo(id nav.ID) tea.Cmd {
	if a.activeScreenID == id {
		return a.resizeLayout()
	}

	next, ok := a.screens[id]
	if !ok {
		return nil
	}

	a.activeScreenID = id
	a.screen = next
	a.keyBar.SetActiveID(id)

	return a.resizeLayout()
}

func (a *App) resizeLayout() tea.Cmd {
	a.keyBar.SetSize(a.width, a.keyBar.Height())

	contentHeight := a.contentHeight()
	nextScreen, screenCmd := a.screen.Update(tea.WindowSizeMsg{
		Width:  a.width,
		Height: contentHeight,
	})
	a.screen = nextScreen
	a.screens[a.activeScreenID] = nextScreen

	return screenCmd
}

func (a *App) contentHeight() int {
	return a.height - a.keyBar.Height()
}

func renderedLineCount(content string) int {
	if content == "" {
		return 0
	}

	return strings.Count(content, "\n") + 1
}
