package repository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"Dzaakk/simple-commerce/internal/catalog/model"
)

var categoryColumns = []string{"id", "parent_id", "name", "slug", "is_active", "created_at", "updated_at"}

func TestCategoryRepositoryCreate(t *testing.T) {
	parentID := int64(10)
	now := time.Date(2026, time.June, 3, 9, 0, 0, 0, time.UTC)
	category := &model.Category{
		ParentID:  &parentID,
		Name:      "Phones",
		Slug:      "phones",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	db, mock := newMockDB(t)
	mock.ExpectQuery(categoryQueryCreate).
		WithArgs(parentID, category.Name, category.Slug, category.IsActive, category.CreatedAt, category.UpdatedAt).
		WillReturnRows(sqlmockRows([]string{"id"}).AddRow(int64(42)))

	got, err := NewCategoryRepository(db).Create(context.Background(), category)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if got != 42 {
		t.Fatalf("id = %d, want 42", got)
	}
}

func TestCategoryRepositoryCreateReturnsWrappedError(t *testing.T) {
	wantErr := errors.New("insert failed")
	db, mock := newMockDB(t)
	mock.ExpectQuery(categoryQueryCreate).
		WithArgs(nil, "Phones", "phones", true, time.Time{}, time.Time{}).
		WillReturnError(wantErr)

	got, err := NewCategoryRepository(db).Create(context.Background(), &model.Category{
		Name:     "Phones",
		Slug:     "phones",
		IsActive: true,
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
	if got != 0 {
		t.Fatalf("id = %d, want 0", got)
	}
}

func TestCategoryRepositoryFindByID(t *testing.T) {
	parentID := int64(10)
	now := time.Date(2026, time.June, 3, 10, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery(categoryQueryFindByID).
		WithArgs(int64(42)).
		WillReturnRows(sqlmockRows(categoryColumns).AddRow(categoryRow(42, &parentID, "Phones", "phones", true, now)...))

	got, err := NewCategoryRepository(db).FindByID(context.Background(), 42)
	if err != nil {
		t.Fatalf("FindByID returned error: %v", err)
	}
	assertCategory(t, got, 42, &parentID, "Phones", "phones", true, now)
}

func TestCategoryRepositoryFindByIDReturnsNotFoundError(t *testing.T) {
	db, mock := newMockDB(t)
	mock.ExpectQuery(categoryQueryFindByID).
		WithArgs(int64(99)).
		WillReturnRows(sqlmockRows(categoryColumns))

	got, err := NewCategoryRepository(db).FindByID(context.Background(), 99)
	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("error = %v, want wrapping sql.ErrNoRows", err)
	}
	if got != nil {
		t.Fatalf("category = %#v, want nil", got)
	}
}

func TestCategoryRepositoryFindAll(t *testing.T) {
	parentID := int64(10)
	now := time.Date(2026, time.June, 3, 11, 0, 0, 0, time.UTC)
	db, mock := newMockDB(t)
	mock.ExpectQuery(categoryQueryFindAll).
		WillReturnRows(sqlmockRows(categoryColumns).
			AddRow(categoryRow(10, nil, "Electronics", "electronics", true, now)...).
			AddRow(categoryRow(42, &parentID, "Phones", "phones", true, now)...))

	got, err := NewCategoryRepository(db).FindAll(context.Background())
	if err != nil {
		t.Fatalf("FindAll returned error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("category count = %d, want 2", len(got))
	}
	assertCategory(t, got[0], 10, nil, "Electronics", "electronics", true, now)
	assertCategory(t, got[1], 42, &parentID, "Phones", "phones", true, now)
}

func TestCategoryRepositoryFindAllReturnsQueryError(t *testing.T) {
	wantErr := errors.New("select failed")
	db, mock := newMockDB(t)
	mock.ExpectQuery(categoryQueryFindAll).WillReturnError(wantErr)

	got, err := NewCategoryRepository(db).FindAll(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want wrapping %v", err, wantErr)
	}
	if got != nil {
		t.Fatalf("categories = %#v, want nil", got)
	}
}

func categoryRow(id int64, parentID *int64, name, slug string, isActive bool, at time.Time) []driver.Value {
	var parent any
	if parentID != nil {
		parent = *parentID
	}
	return []driver.Value{id, parent, name, slug, isActive, at, at}
}

func assertCategory(t *testing.T, got *model.Category, id int64, parentID *int64, name, slug string, isActive bool, at time.Time) {
	t.Helper()

	if got == nil {
		t.Fatal("category is nil")
	}
	if got.ID != id || got.Name != name || got.Slug != slug || got.IsActive != isActive ||
		!got.CreatedAt.Equal(at) || !got.UpdatedAt.Equal(at) {
		t.Fatalf("category = %#v", got)
	}
	if parentID == nil {
		if got.ParentID != nil {
			t.Fatalf("parent id = %v, want nil", *got.ParentID)
		}
		return
	}
	if got.ParentID == nil || *got.ParentID != *parentID {
		t.Fatalf("parent id = %v, want %d", got.ParentID, *parentID)
	}
}
