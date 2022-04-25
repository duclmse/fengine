package sql

import (
	"context"
	"database/sql"
	"github.com/duclmse/fengine/pkg/logger"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

// Database provides a database interface
type Database interface {
	Connection() *sqlx.DB

	NamedExecContext(context.Context, string, any) (sql.Result, error)

	QueryRowxContext(context.Context, string, ...any) *sqlx.Row

	NamedQueryContext(context.Context, string, any) (*sqlx.Rows, error)

	GetContext(context.Context, any, string, ...any) error

	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)

	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
}

// NewDatabase creates a DeviceDatabase instance
func NewDatabase(DB *sqlx.DB) Database {
	return &database{
		DB: DB,
	}
}

type database struct {
	DB  *sqlx.DB
	Log logger.Logger
}

func (db database) Connection() *sqlx.DB {
	return db.DB
}

func (db database) Logger() logger.Logger {
	return db.Log
}

func (db database) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	addSpanTags(ctx, query)
	return db.DB.ExecContext(ctx, query, args...)
}

func (db database) NamedExec(ctx context.Context, query string, args any) (sql.Result, error) {
	return db.DB.NamedExecContext(ctx, query, args)
}

func (db database) QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	addSpanTags(ctx, query)
	return db.DB.QueryxContext(ctx, query, args...)
}

func (db database) NamedExecContext(ctx context.Context, query string, args any) (sql.Result, error) {
	addSpanTags(ctx, query)
	return db.DB.NamedExecContext(ctx, query, args)
}

func (db database) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	addSpanTags(ctx, query)
	return db.DB.QueryRowxContext(ctx, query, args...)
}

func (db database) NamedQueryContext(ctx context.Context, query string, args any) (*sqlx.Rows, error) {
	addSpanTags(ctx, query)
	return db.DB.NamedQueryContext(ctx, query, args)
}

func (db database) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	addSpanTags(ctx, query)
	return db.DB.GetContext(ctx, dest, query, args...)
}

func (db database) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("span.kind", "client")
		span.SetTag("peer.service", "postgres")
		span.SetTag("db.type", "sql")
	}
	return db.DB.BeginTxx(ctx, opts)
}

func addSpanTags(ctx context.Context, query string) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag("sql.statement", query)
		span.SetTag("span.kind", "client")
		span.SetTag("peer.service", "postgres")
		span.SetTag("db.type", "sql")
	}
}
