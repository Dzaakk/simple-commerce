package repository

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTransactionRepositoryGenerateTransactionNumberUsesAtomicCounter(t *testing.T) {
	db, mock := newTransactionMockDB(t)
	mock.ExpectQuery("INSERT INTO public.business_number_counters").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(int64(34)))

	got, err := NewTransactionRepository(db).GenerateTransactionNumber(context.Background())
	if err != nil {
		t.Fatalf("GenerateTransactionNumber returned error: %v", err)
	}

	today := time.Now().Format("20060102")
	wantPrefix := "TRX-" + today + "-0034"
	if got != wantPrefix {
		t.Fatalf("transaction number = %q, want %q", got, wantPrefix)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
	if strings.Contains(transactionQueryNextNumber, "COUNT(*)") {
		t.Fatal("transaction number generator must not use COUNT(*)")
	}
}

func newTransactionMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})

	return db, mock
}
