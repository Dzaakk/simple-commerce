package transactor

import (
	"context"
	"database/sql"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Transactor interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type SQLTransactor struct {
	db *sql.DB
}

type executorContextKey struct{}

func New(db *sql.DB) *SQLTransactor {
	return &SQLTransactor{db: db}
}

func (t *SQLTransactor) WithinTx(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	return fn(WithExecutor(ctx, tx))
}

func WithExecutor(ctx context.Context, executor Executor) context.Context {
	return context.WithValue(ctx, executorContextKey{}, executor)
}

func ExecutorFrom(ctx context.Context, fallback Executor) Executor {
	if executor, ok := ctx.Value(executorContextKey{}).(Executor); ok && executor != nil {
		return executor
	}
	return fallback
}
