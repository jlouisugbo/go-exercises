package registration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testApp struct {
	handler *RegistrationHandler
	users   *MemoryUserStore
	emails  *MemoryEmailSender
}

type registrationResult struct {
	status     int
	user       User
	userCount  int
	emailCount int
}

// These helpers are expected to change as you refactor. Update them to match
// your new API. The tests below should not need broad rewrites.
func newTestApp() *testApp {
	users := NewMemoryUserStore()
	emails := &MemoryEmailSender{}
	return &testApp{
		handler: NewRegistrationHandler(users, emails),
		users:   users,
		emails:  emails,
	}
}

func (a *testApp) register(method, body string) registrationResult {
	request := httptest.NewRequest(method, "/registrations", bytes.NewBufferString(body))
	response := httptest.NewRecorder()
	a.handler.ServeHTTP(response, request)

	var user User
	_ = json.Unmarshal(response.Body.Bytes(), &user)

	return registrationResult{
		status:     response.Code,
		user:       user,
		userCount:  a.users.Count(),
		emailCount: a.emails.Count(),
	}
}

func TestServiceRegisterWithoutHTTP(t *testing.T) {
	users := NewMemoryUserStore()
	emails := &MemoryEmailSender{}
	service := NewService(users, emails)

	user, err := service.Register("  JOEL@Example.com ", ProPlan)
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if user.Email != "joel@example.com" {
		t.Fatalf("email = %q, want %q", user.Email, "joel@example.com")
	}
	if user.Plan != ProPlan {
		t.Fatalf("plan = %q, want %q", user.Plan, ProPlan)
	}
	if users.Count() != 1 {
		t.Fatalf("user count = %d, want 1", users.Count())
	}
	if emails.Count() != 1 {
		t.Fatalf("email count = %d, want 1", emails.Count())
	}
}

func TestRegistrationHandler(t *testing.T) {
	t.Run("request handling", func(t *testing.T) {
		t.Run("rejects methods other than POST", func(t *testing.T) {
			app := newTestApp()
			result := app.register(http.MethodGet, "")
			if result.status != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want %d", result.status, http.StatusMethodNotAllowed)
			}
		})

		t.Run("rejects malformed JSON", func(t *testing.T) {
			app := newTestApp()
			result := app.register(http.MethodPost, `{`)
			if result.status != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", result.status, http.StatusBadRequest)
			}
		})

		t.Run("rejects unknown fields", func(t *testing.T) {
			app := newTestApp()
			result := app.register(http.MethodPost, `{"email":"joel@example.com","plan":"free","admin":true}`)
			if result.status != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", result.status, http.StatusBadRequest)
			}
		})
	})

	t.Run("registration rules", func(t *testing.T) {
		t.Run("rejects invalid email", func(t *testing.T) {
			app := newTestApp()
			result := app.register(http.MethodPost, `{"email":"missing-at-sign","plan":"free"}`)
			if result.status != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", result.status, http.StatusBadRequest)
			}
		})

		t.Run("rejects invalid plan", func(t *testing.T) {
			app := newTestApp()
			result := app.register(http.MethodPost, `{"email":"joel@example.com","plan":"enterprise"}`)
			if result.status != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", result.status, http.StatusBadRequest)
			}
		})

		t.Run("rejects duplicate email", func(t *testing.T) {
			app := newTestApp()
			_ = app.register(http.MethodPost, `{"email":"joel@example.com","plan":"free"}`)
			result := app.register(http.MethodPost, `{"email":"JOEL@example.com","plan":"pro"}`)
			if result.status != http.StatusConflict {
				t.Fatalf("status = %d, want %d", result.status, http.StatusConflict)
			}
		})

		t.Run("normalizes email in response", func(t *testing.T) {
			app := newTestApp()
			result := app.register(http.MethodPost, `{"email":"  JOEL@Example.com ","plan":"pro"}`)
			if result.user.Email != "joel@example.com" {
				t.Fatalf("email = %q, want %q", result.user.Email, "joel@example.com")
			}
		})
	})

	t.Run("successful registration", func(t *testing.T) {
		t.Run("returns created", func(t *testing.T) {
			app := newTestApp()
			result := app.register(http.MethodPost, `{"email":"joel@example.com","plan":"team"}`)
			if result.status != http.StatusCreated {
				t.Fatalf("status = %d, want %d", result.status, http.StatusCreated)
			}
		})

		t.Run("stores one user", func(t *testing.T) {
			app := newTestApp()
			result := app.register(http.MethodPost, `{"email":"joel@example.com","plan":"team"}`)
			if result.userCount != 1 {
				t.Fatalf("user count = %d, want 1", result.userCount)
			}
		})

		t.Run("sends one welcome email", func(t *testing.T) {
			app := newTestApp()
			result := app.register(http.MethodPost, `{"email":"joel@example.com","plan":"team"}`)
			if result.emailCount != 1 {
				t.Fatalf("email count = %d, want 1", result.emailCount)
			}
		})
	})

	t.Run("dependency failures", func(t *testing.T) {
		t.Run("returns server error when saving fails", func(t *testing.T) {
			app := newTestApp()
			app.users.SaveErr = ErrUnavailable
			result := app.register(http.MethodPost, `{"email":"joel@example.com","plan":"free"}`)
			if result.status != http.StatusInternalServerError {
				t.Fatalf("status = %d, want %d", result.status, http.StatusInternalServerError)
			}
		})

		t.Run("does not email when saving fails", func(t *testing.T) {
			app := newTestApp()
			app.users.SaveErr = ErrUnavailable
			result := app.register(http.MethodPost, `{"email":"joel@example.com","plan":"free"}`)
			if result.emailCount != 0 {
				t.Fatalf("email count = %d, want 0", result.emailCount)
			}
		})

		t.Run("returns server error when email fails", func(t *testing.T) {
			app := newTestApp()
			app.emails.SendErr = ErrUnavailable
			result := app.register(http.MethodPost, `{"email":"joel@example.com","plan":"free"}`)
			if result.status != http.StatusInternalServerError {
				t.Fatalf("status = %d, want %d", result.status, http.StatusInternalServerError)
			}
		})
	})
}
