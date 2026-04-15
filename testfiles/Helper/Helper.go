package helper

import (
	"log/slog"
	"os"
	"sync"
)

var (
	defaultLogger *slog.Logger
	once          sync.Once
)

// Init initializes the global logger (call once in main)
func Init(filename string, level slog.Level) error {
	var err error

	once.Do(func() {
		file, e := os.OpenFile(
			filename,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0666,
		)
		if e != nil {
			err = e
			return
		}

		opts := &slog.HandlerOptions{
			Level:     level,
			AddSource: true, // adds file + line
		}

		handler := slog.NewJSONHandler(file, opts)
		defaultLogger = slog.New(handler)

		// Make it default for slog package usage
		slog.SetDefault(defaultLogger)
	})

	return err
}

// L returns the global logger
func L() *slog.Logger {
	if defaultLogger == nil {
		panic("logger not initialized")
	}
	return defaultLogger
}
