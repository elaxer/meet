package database

import (
	"context"
	"database/sql"
	"log/slog"
	"meet/internal/pkg/app/slogger"
)

type connectionLogging struct {
	Connection
}

func NewConnectionLogging(connection Connection) Connection {
	return &connectionLogging{connection}
}

func (cl *connectionLogging) Exec(query string, args ...any) (sql.Result, error) {
	slog.Log(context.Background(), slogger.LevelSQL, "Выполнение запроса метода Exec", "query", query, "args", args)

	return cl.Connection.Exec(query, args...)
}

func (cl *connectionLogging) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	slog.Log(ctx, slogger.LevelSQL, "Выполнение запроса метода ExecContext", "query", query, "args", args)

	return cl.Connection.ExecContext(ctx, query, args...)
}

func (cl *connectionLogging) Query(query string, args ...any) (*sql.Rows, error) {
	slog.Log(context.Background(), slogger.LevelSQL, "Выполнение запроса метода Query", "query", query, "args", args)

	return cl.Connection.Query(query, args...)
}

func (cl *connectionLogging) QueryRow(query string, args ...any) *sql.Row {
	slog.Log(context.Background(), slogger.LevelSQL, "Выполнение запроса метода QueryRow", "query", query, "args", args)

	return cl.Connection.QueryRow(query, args...)
}
