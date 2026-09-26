package util

import "testing"

func TestAccessTokenRoundTrip(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")
	raw, err := GenerateAccessToken("customer-1", "customer@example.com")
	if err != nil {
		t.Fatal(err)
	}
	claims, err := ParseAccessToken(raw)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "customer-1" || claims.Subject != "customer-1" || claims.Email != "customer@example.com" {
		t.Fatalf("claims = %#v", claims)
	}
}

func TestAccessTokenRequiresSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")
	if _, err := GenerateAccessToken("customer-1", "customer@example.com"); err == nil {
		t.Fatal("GenerateAccessToken() error = nil")
	}
}
