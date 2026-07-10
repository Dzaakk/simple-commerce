package repository

import (
	"context"
	"database/sql"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestOrderRepositoryGenerateOrderNumberUsesAtomicCounter(t *testing.T) {
	db, mock := newOrderMockDB(t)
	mock.ExpectQuery("INSERT INTO public.business_number_counters").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"value"}).AddRow(int64(12)))

	got, err := NewOrderRepository(db).GenerateOrderNumber(context.Background())
	if err != nil {
		t.Fatalf("GenerateOrderNumber returned error: %v", err)
	}

	today := time.Now().Format("20060102")
	wantPrefix := "ORD-" + today + "-0012"
	if got != wantPrefix {
		t.Fatalf("order number = %q, want %q", got, wantPrefix)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
	if strings.Contains(orderQueryNextNumber, "COUNT(*)") {
		t.Fatal("order number generator must not use COUNT(*)")
	}
}

func newOrderMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
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
