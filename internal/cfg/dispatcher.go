package cfg

import (
	vein "github.com/igorrius/go-vein"
	"github.com/samber/do/v2"
)

// ProvideDispatcher returns a do provider that creates the application-wide go-vein Dispatcher.
func ProvideDispatcher() do.Provider[*vein.Dispatcher] {
	return func(_ do.Injector) (*vein.Dispatcher, error) {
		return &vein.Dispatcher{}, nil
	}
}
