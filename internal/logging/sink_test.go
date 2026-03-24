package logging_test

import (
	"os"
	"testing"
	"time"

	"github.com/igorrius/go-app-tui-template/internal/config"
	"github.com/igorrius/go-app-tui-template/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSinkWriter_Stdout(t *testing.T) {
	cfg := config.LoggingConfig{Sink: "stdout"}
	w, closer, err := logging.NewSinkWriterForTest(cfg)
	require.NoError(t, err)
	assert.Equal(t, os.Stdout, w)
	assert.Nil(t, closer)
}

func TestNewSinkWriter_Stderr(t *testing.T) {
	cfg := config.LoggingConfig{Sink: "stderr"}
	w, closer, err := logging.NewSinkWriterForTest(cfg)
	require.NoError(t, err)
	assert.Equal(t, os.Stderr, w)
	assert.Nil(t, closer)
}

func TestNewSinkWriter_UnknownSink(t *testing.T) {
	cfg := config.LoggingConfig{Sink: "nats"}
	_, _, err := logging.NewSinkWriterForTest(cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "nats")
}

func TestNewSinkWriter_FileCreatesDir(t *testing.T) {
	dir := t.TempDir() + "/nested/logs"
	cfg := config.LoggingConfig{
		Sink:             "file",
		Dir:              dir,
		RotationInterval: time.Hour,
	}
	w, closer, err := logging.NewSinkWriterForTest(cfg)
	require.NoError(t, err)
	require.NotNil(t, w)
	require.NotNil(t, closer)
	defer closer.Close() // nolint:errcheck

	_, statErr := os.Stat(dir)
	assert.NoError(t, statErr, "log directory should have been created")
}
