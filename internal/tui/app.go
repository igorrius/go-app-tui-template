package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/igorrius/go-app-tui-template/internal/config"
	"github.com/igorrius/go-app-tui-template/internal/logging"
	"github.com/igorrius/go-app-tui-template/internal/tui/component"
	"github.com/igorrius/go-app-tui-template/internal/tui/screen"
	vein "github.com/igorrius/go-vein"
)

// App is the root bubbletea model managing active screen, terminal size, and key bar.
type App struct {
	cfg            *config.AppConfig
	dispatcher     *vein.Dispatcher
	logCh          <-chan logging.LogEvent
	activeScreenID screen.ID
	screen         tea.Model
	keyBar         *component.KeyBar
	width          int
	height         int
}

// NewApp creates a new root TUI application model.
func NewApp(cfg *config.AppConfig, dispatcher *vein.Dispatcher) *App {
	sub := vein.SubscribeTo[logging.LogEvent](dispatcher)
	logCh := sub.OnC()

	initialScreen, err := screen.New(screen.DashboardID, logCh)
	if err != nil {
		panic(err)
	}

	return &App{
		cfg:            cfg,
		dispatcher:     dispatcher,
		logCh:          logCh,
		activeScreenID: screen.DashboardID,
		screen:         initialScreen,
		keyBar:         component.NewKeyBar(),
	}
}

func (a *App) Init() tea.Cmd {
	return tea.Batch(a.screen.Init(), a.keyBar.Init())
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, a.resizeLayout()

	case tea.KeyMsg:
		switch msg.String() {
		case "f10", "ctrl+c":
			return a, tea.Quit
		case "f1":
			return a, a.routeTo(screen.DashboardID)
		case "f9":
			return a, a.routeTo(screen.LogViewID)
		}
	}

	nextScreen, screenCmd := a.screen.Update(msg)
	a.screen = nextScreen

	nextKeyBar, keyBarCmd := a.keyBar.Update(msg)
	if kb, ok := nextKeyBar.(*component.KeyBar); ok {
		a.keyBar = kb
	}

	return a, tea.Batch(screenCmd, keyBarCmd)
}

func (a *App) View() tea.View {
	v := tea.NewView("App")
	v.AltScreen = true

	if a.width == 0 || a.height == 0 {
		return v
	}

	screenView := a.screen.View()
	keyBarView := a.keyBar.View()
	screenLines := renderedLineCount(screenView.Content)
	contentHeight := a.contentHeight()

	// Pad screen content to push key bar to the bottom
	var b strings.Builder
	b.WriteString(screenView.Content)
	padLines := contentHeight - screenLines
	if screenLines > 0 {
		padLines++
	}
	for i := 0; i < padLines; i++ {
		b.WriteString("\n")
	}
	b.WriteString(keyBarView.Content)

	v.SetContent(b.String())
	return v
}

func (a *App) routeTo(id screen.ID) tea.Cmd {
	if a.activeScreenID == id {
		return a.resizeLayout()
	}

	next, err := screen.New(id, a.logCh)
	if err != nil {
		return nil
	}

	a.activeScreenID = id
	a.screen = next
	a.keyBar.SetActiveID(id)

	return tea.Batch(next.Init(), a.resizeLayout())
}

func (a *App) resizeLayout() tea.Cmd {
	a.keyBar.SetSize(a.width, a.keyBar.Height())

	nextScreen, screenCmd := a.screen.Update(tea.WindowSizeMsg{
		Width:  a.width,
		Height: a.contentHeight(),
	})
	a.screen = nextScreen

	return screenCmd
}

func (a *App) contentHeight() int {
	keyBarHeight := a.keyBar.Height()
	if a.height <= keyBarHeight {
		return 0
	}

	return a.height - keyBarHeight
}

func renderedLineCount(content string) int {
	if content == "" {
		return 0
	}

	return strings.Count(content, "\n") + 1
}
