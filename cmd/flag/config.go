package flag

import "github.com/urfave/cli/v3"

const (
	// ConfigFilePathFlagName is the name of the CLI flag for specifying the config file path.
	ConfigFilePathFlagName = "config"
	// ConsoleModeFlag is the name of the CLI flag for running in console mode.
	ConsoleModeFlag = "console"
)

var ConfigFilePath cli.Flag = &cli.StringFlag{
	Name:    ConfigFilePathFlagName,
	Usage:   "path to config file",
	Sources: cli.EnvVars("APP_CONFIG"),
	Value:   "app-config.dist.yaml",
}

var ConsoleMode cli.Flag = &cli.BoolFlag{
	Name:    ConsoleModeFlag,
	Usage:   "run in console mode (no TUI); all logs are written to stdout",
	Sources: cli.EnvVars("APP_CONSOLE"),
}
