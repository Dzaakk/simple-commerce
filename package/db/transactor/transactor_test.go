package transactor

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSQLTransactorWithinTxCommitsOnSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	err = New(db).WithinTx(context.Background(), func(ctx context.Context) error {
		if ExecutorFrom(ctx, nil) == nil {
			t.Fatal("transaction executor must be available in context")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithinTx returned error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestSQLTransactorWithinTxRollsBackOnError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	defer db.Close()

	wantErr := errors.New("write failed")
	mock.ExpectBegin()
	mock.ExpectRollback()

	err = New(db).WithinTx(context.Background(), func(context.Context) error {
		return wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
