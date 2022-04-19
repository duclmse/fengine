package sql

import (
	"context"
	"database/sql"
	"github.com/duclmse/fengine/pkg/logger"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
)

type database struct {
	DB  *sqlx.DB
	Log logger.Logger
}

// Database provides a database interface
type Database interface {
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)

	QueryRowxContext(context.Context, string, ...interface{}) *sqlx.Row

	NamedQueryContext(context.Context, string, interface{}) (*sqlx.Rows, error)

	GetContext(context.Context, interface{}, string, ...interface{}) error

	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
}

// NewDatabase creates a DeviceDatabase instance
func NewDatabase(DB *sqlx.DB) Database {
	return &database{
		DB: DB,
	}
}

func (db database) Logger() logger.Logger {
	return db.Log
}

func (db database) NamedExecContext(ctx context.Context, query string, args interface{}) (sql.Result, error) {
	addSpanTags(ctx, query)
	return db.DB.NamedExecContext(ctx, query, args)
}

func (db database) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	addSpanTags(ctx, query)
	return db.DB.QueryRowxContext(ctx, query, args...)
}

func (db database) NamedQueryContext(ctx context.Context, query string, args interface{}) (*sqlx.Rows, error) {
	addSpanTags(ctx, query)
	return db.DB.NamedQueryContext(ctx, query, args)
}

func (db database) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
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
