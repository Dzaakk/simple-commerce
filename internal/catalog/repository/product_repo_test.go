package repository

import (
	"reflect"
	"testing"

	"Dzaakk/simple-commerce/internal/catalog/dto"
)

func TestBuildProductQueryUsesStableOrderingAndFilters(t *testing.T) {
	categoryID := int64(2)
	minPrice := 1000.0
	name := "phone"
	query, args := buildProductQuery(dto.ProductQuery{
		CategoryID: &categoryID, MinPrice: &minPrice, Name: &name, Limit: 20, SortBy: "price_asc",
	})
	wantQuery := "SELECT " + productColumns + " FROM products WHERE is_active = true AND category_id = $1 AND price >= $2 AND name ILIKE $3 ORDER BY price ASC, id ASC LIMIT $4"
	if query != wantQuery {
		t.Fatalf("query = %q\nwant  = %q", query, wantQuery)
	}
	if want := []any{int64(2), 1000.0, "%phone%", 20}; !reflect.DeepEqual(args, want) {
		t.Fatalf("args = %#v, want %#v", args, want)
	}
}

func TestSplitCursorRejectsInvalidID(t *testing.T) {
	if _, _, ok := splitCursor("value|not-an-id"); ok {
		t.Fatal("splitCursor() accepted invalid id")
	}
}
