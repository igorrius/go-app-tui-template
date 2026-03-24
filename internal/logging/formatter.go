package logging

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// TextFormatter renders a LogEvent as a single human-readable line in the format:
//
//	<Time> <LEVEL> <MODULE> [key=value ...]
//
// Time uses the time.DateTime layout. Level and Module are always uppercased.
// Module is extracted from the "module" attribute (default: "APP") and is excluded
// from the trailing key=value pairs.
type TextFormatter struct{}

// Format renders e as a single human-readable log line.
func (f TextFormatter) Format(e LogEvent) string {
	module := "APP"
	var remaining []slog.Attr
	for _, a := range e.Attrs {
		if a.Key == "module" {
			module = strings.ToUpper(a.Value.String())
		} else {
			remaining = append(remaining, a)
		}
	}

	var b strings.Builder
	b.WriteByte('[')
	b.WriteString(e.Time.Format(time.DateTime))
	b.WriteString("] ")
	b.WriteString(strings.ToUpper(e.Level.String()))
	b.WriteByte(' ')
	b.WriteString(module)
	b.WriteByte(' ')
	b.WriteString(e.Message)

	for _, a := range remaining {
		fmt.Fprintf(&b, " %s=%s", a.Key, a.Value.String())
	}

	return b.String()
}

// ColorTextFormatter renders a LogEvent as a single human-readable line with ANSI colors.
// It uses the same layout as TextFormatter but applies styling to various fields.
type ColorTextFormatter struct{}

// Format renders e as a single human-readable log line with ANSI colors.
func (f ColorTextFormatter) Format(e LogEvent) string {
	timestampStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
	moduleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
	attrKeyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))

	var levelStyle lipgloss.Style
	var levelStr string
	switch e.Level {
	case slog.LevelDebug:
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")) // Blue
		levelStr = "DEBUG"
	case slog.LevelInfo:
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
		levelStr = "INFO "
	case slog.LevelWarn:
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Yellow
		levelStr = "WARN "
	case slog.LevelError:
		levelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red
		levelStr = "ERROR"
	default:
		levelStyle = lipgloss.NewStyle()
		levelStr = fmt.Sprintf("%-5s", strings.ToUpper(e.Level.String()))
	}

	module := "APP"
	var remaining []slog.Attr
	for _, a := range e.Attrs {
		if a.Key == "module" {
			module = strings.ToUpper(a.Value.String())
		} else {
			remaining = append(remaining, a)
		}
	}

	var b strings.Builder
	b.WriteString(timestampStyle.Render("[" + e.Time.Format(time.DateTime) + "]"))
	b.WriteByte(' ')
	b.WriteString(levelStyle.Render(levelStr))
	b.WriteByte(' ')
	b.WriteString(moduleStyle.Render(module))
	b.WriteByte(' ')
	b.WriteString(e.Message)

	for _, a := range remaining {
		fmt.Fprintf(&b, " %s=%s", attrKeyStyle.Render(a.Key), a.Value.String())
	}

	return b.String()
}
