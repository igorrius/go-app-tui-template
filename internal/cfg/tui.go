package cfg

import (
	"github.com/igorrius/go-app-tui-template/internal/config"
	"github.com/igorrius/go-app-tui-template/internal/tui"
	vein "github.com/igorrius/go-vein"
	"github.com/samber/do/v2"
)

// ProvideTUIApp returns a do provider that creates the root TUI application model.
func ProvideTUIApp() do.Provider[*tui.App] {
	return func(i do.Injector) (*tui.App, error) {
		cfg, err := do.Invoke[*config.AppConfig](i)
		if err != nil {
			return nil, err
		}

		dispatcher, err := do.Invoke[*vein.Dispatcher](i)
		if err != nil {
			return nil, err
		}

		return tui.NewApp(cfg, dispatcher), nil
	}
}
