package cfg

import (
	"github.com/igorrius/go-app-tui-template/cmd/flag"
	"github.com/igorrius/go-app-tui-template/internal/config"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
)

// ProvideConfig returns a do provider that loads and returns the application config.
func ProvideConfig(cmd *cli.Command) do.Provider[*config.AppConfig] {
	return func(i do.Injector) (*config.AppConfig, error) {
		return config.Load(cmd.String(flag.ConfigFilePathFlagName))
	}
}
