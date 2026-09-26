package main

import (
	"os"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSchemaOnlyContainsRevampDomains(t *testing.T) {
	contents, err := os.ReadFile("db.sql")
	if err != nil {
		t.Fatal(err)
	}
	schema := strings.ToLower(string(contents))
	for _, table := range []string{"users", "refresh_tokens", "categories", "products"} {
		if !strings.Contains(schema, "create table "+table) {
			t.Errorf("schema is missing %s", table)
		}
	}
	for _, removed := range []string{"sellers", "carts", "cart_items", "inventories", "orders", "transactions", "activation_codes"} {
		if strings.Contains(schema, "create table "+removed) {
			t.Errorf("removed table %s is still present", removed)
		}
	}
}

func TestRouterOnlyExposesRevampScope(t *testing.T) {
	want := map[string]bool{
		"GET /healthz": true, "GET /readyz": true,
		"POST /api/v1/auth/customer/register": true,
		"POST /api/v1/auth/login":             true,
		"POST /api/v1/auth/refresh-token":     true,
		"POST /api/v1/auth/logout":            true,
		"GET /api/v1/customer/me":             true,
		"GET /api/v1/category":                true,
		"GET /api/v1/product":                 true,
		"GET /api/v1/product/:id":             true,
		"GET /api/v2/product":                 true,
		"GET /api/v2/product/:id":             true,
	}
	routes := newRouter(nil, nil).Routes()
	if len(routes) != len(want) {
		t.Fatalf("route count = %d, want %d: %#v", len(routes), len(want), routes)
	}
	for _, route := range routes {
		key := route.Method + " " + route.Path
		if !want[key] {
			t.Errorf("unexpected route %s", key)
		}
	}
}

func TestOpenAPIOnlyExposesRevampScope(t *testing.T) {
	contents, err := os.ReadFile("api.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var contract struct {
		Paths map[string]any `yaml:"paths"`
	}
	if err := yaml.Unmarshal(contents, &contract); err != nil {
		t.Fatalf("parse api.yaml: %v", err)
	}

	want := map[string]bool{
		"/healthz": true, "/readyz": true,
		"/api/v1/auth/customer/register": true,
		"/api/v1/auth/login":             true,
		"/api/v1/auth/refresh-token":     true,
		"/api/v1/auth/logout":            true,
		"/api/v1/customer/me":            true,
		"/api/v1/category":               true,
		"/api/v1/product":                true,
		"/api/v1/product/{id}":           true,
		"/api/v2/product":                true,
		"/api/v2/product/{id}":           true,
	}
	if len(contract.Paths) != len(want) {
		t.Fatalf("path count = %d, want %d: %#v", len(contract.Paths), len(want), contract.Paths)
	}
	for path := range want {
		if _, ok := contract.Paths[path]; !ok {
			t.Errorf("missing path %s", path)
		}
	}
}
