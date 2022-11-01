package log

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// Logger implements logging methods, adding tracing fields whenever they are provided.
type Logger struct {
	zerolog.Logger
}

// NewLogger instantiates a *Logger.
func NewLogger() *Logger {
	return newLogger(os.Stderr)
}

func newLogger(w io.Writer) *Logger {
	return &Logger{
		Logger: zerolog.New(w).With().Timestamp().Logger(),
	}
}
