package service

import (
	"log/slog"
	"testing"
	"time"

	"github.com/givetrack/givetrack/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

func newTestAuthService() *AuthService {
	return NewAuthService(nil, nil, "unit-test-secret-which-is-long-enough", 1, slog.Default())
}

func TestGenerateAndParseToken(t *testing.T) {
	svc := newTestAuthService()
	tests := []struct {
		name string
		user model.User
	}{
		{name: "admin", user: model.User{ID: 1, Username: "admin", Role: "admin"}},
		{name: "org", user: model.User{ID: 2, Username: "careorg", Role: "org"}},
		{name: "donor", user: model.User{ID: 3, Username: "donor1", Role: "user"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := svc.GenerateToken(&tt.user)
			if err != nil {
				t.Fatalf("GenerateToken: %v", err)
			}
			claims, err := svc.ParseToken(token)
			if err != nil {
				t.Fatalf("ParseToken: %v", err)
			}
			if claims.UserID != tt.user.ID || claims.Role != tt.user.Role || claims.Subject != tt.user.Username {
				t.Fatalf("claims mismatch: %+v", claims)
			}
		})
	}
}

func TestParseTokenRejectsInvalid(t *testing.T) {
	svc := newTestAuthService()
	for _, tok := range []string{"", "not-a-jwt", "eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJhZG1pbiJ9.invalid"} {
		if _, err := svc.ParseToken(tok); err == nil {
			t.Fatalf("expected error for %q", tok)
		}
	}
}

func TestParseTokenRejectsExpired(t *testing.T) {
	svc := NewAuthService(nil, nil, "unit-test-secret-which-is-long-enough", 0, slog.Default())
	now := time.Now()
	claims := Claims{UserID: 1, Role: "user", RegisteredClaims: jwt.RegisteredClaims{Subject: "donor1", ExpiresAt: jwt.NewNumericDate(now.Add(-time.Hour))}}
	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(svc.jwtSecret)
	if _, err := svc.ParseToken(token); err == nil {
		t.Fatal("expected expired token error")
	}
}
