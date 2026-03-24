package logging

import (
	"context"
	"log/slog"

	vein "github.com/igorrius/go-vein"
)

// VeinHandler is a slog.Handler that publishes a LogEvent to the go-vein bus on each Handle call.
type VeinHandler struct {
	level  slog.Level
	bus    *vein.Dispatcher
	attrs  []slog.Attr
	groups []string
}

// NewVeinHandler creates a VeinHandler that filters below minLevel and publishes to bus.
func NewVeinHandler(minLevel slog.Level, bus *vein.Dispatcher) *VeinHandler {
	return &VeinHandler{level: minLevel, bus: bus}
}

// Enabled reports whether the handler handles records at the given level.
func (h *VeinHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle converts r into a LogEvent and publishes it to the bus.
func (h *VeinHandler) Handle(_ context.Context, r slog.Record) error {
	attrs := make([]slog.Attr, 0, len(h.attrs)+r.NumAttrs())
	attrs = append(attrs, h.attrs...)

	var recordAttrs []slog.Attr
	r.Attrs(func(a slog.Attr) bool {
		recordAttrs = append(recordAttrs, a)
		return true
	})
	attrs = append(attrs, wrapAttrsInGroups(recordAttrs, h.groups)...)

	vein.PublishTo(h.bus, LogEvent{
		Time:    r.Time,
		Level:   r.Level,
		Message: r.Message,
		Attrs:   attrs,
	})
	return nil
}

// WithAttrs returns a new VeinHandler whose log records include attrs.
func (h *VeinHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	newAttrs := make([]slog.Attr, len(h.attrs), len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	newAttrs = append(newAttrs, wrapAttrsInGroups(attrs, h.groups)...)
	return &VeinHandler{
		level:  h.level,
		bus:    h.bus,
		attrs:  newAttrs,
		groups: h.groups,
	}
}

// WithGroup returns a new VeinHandler that nests subsequent attributes under name.
func (h *VeinHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	newGroups := make([]string, len(h.groups)+1)
	copy(newGroups, h.groups)
	newGroups[len(h.groups)] = name
	return &VeinHandler{
		level:  h.level,
		bus:    h.bus,
		attrs:  h.attrs,
		groups: newGroups,
	}
}

// wrapAttrsInGroups nests attrs under the group hierarchy (outermost first).
func wrapAttrsInGroups(attrs []slog.Attr, groups []string) []slog.Attr {
	if len(groups) == 0 || len(attrs) == 0 {
		return attrs
	}
	result := attrs
	for i := len(groups) - 1; i >= 0; i-- {
		args := make([]any, len(result))
		for j, a := range result {
			args[j] = a
		}
		result = []slog.Attr{slog.Group(groups[i], args...)}
	}
	return result
}
