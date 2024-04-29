package helper

import (
	"io"
	"log/slog"
	"os"
)

func ConfigureSlogger(debug bool, writers ...io.Writer) {
	writers = append(writers, os.Stdout)
	multiW := io.MultiWriter(writers...)

	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(multiW, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
