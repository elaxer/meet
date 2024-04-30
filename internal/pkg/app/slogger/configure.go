package slogger

import (
	"io"
	"log/slog"
	"os"
)

var logLevelNames = map[slog.Leveler]string{
	LevelSQL: "SQL",
}

func Configure(debug bool, writers ...io.Writer) {
	writers = append(writers, os.Stdout)
	multiW := io.MultiWriter(writers...)

	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(multiW, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := logLevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
