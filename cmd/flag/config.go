package flag

import "github.com/urfave/cli/v3"

const (
	// ConfigFilePathFlagName is the name of the CLI flag for specifying the config file path.
	ConfigFilePathFlagName = "config"
)

var ConfigFilePath cli.Flag = &cli.StringFlag{
	Name:    ConfigFilePathFlagName,
	Usage:   "path to config file",
	Sources: cli.EnvVars("GO_APP_TUI_TEMPLATE_CONFIG"),
	Value:   "app-config.dist.yaml",
}
