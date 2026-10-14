package payroll

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// spyHandler records whether it was invoked, standing in for whatever
// handler sits behind a middleware under test.
type spyHandler struct {
	called bool
}

func (s *spyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.called = true
	w.WriteHeader(http.StatusOK)
}

// These helpers are expected to change as you refactor -- update them to match your new API.
func callRequireRole(requiredRole, callerRole string) (status int, downstream *spyHandler) {
	downstream = &spyHandler{}
	handler := RequireRole(requiredRole, downstream)

	req := httptest.NewRequest(http.MethodPost, "/adjust", nil)
	if callerRole != "" {
		req.Header.Set("X-Employee-Role", callerRole)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec.Code, downstream
}

func callPayrollHandler(employeeID, callerRole string) (status int, body string) {
	handler := NewPayrollHandler()

	req := httptest.NewRequest(http.MethodPost, "/adjust?employee_id="+employeeID, nil)
	if callerRole != "" {
		req.Header.Set("X-Employee-Role", callerRole)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	return rec.Code, rec.Body.String()
}

// End of helpers -- the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestRequireRole(t *testing.T) {
	t.Run("matching role", func(t *testing.T) {
		t.Run("lets the request through", func(t *testing.T) {
			_, downstream := callRequireRole("admin", "admin")
			if !downstream.called {
				t.Fatal("expected downstream handler to run for a matching role")
			}
		})

		t.Run("returns the downstream status", func(t *testing.T) {
			status, _ := callRequireRole("admin", "admin")
			if status != http.StatusOK {
				t.Fatalf("status = %d, want %d", status, http.StatusOK)
			}
		})
	})

	t.Run("mismatched role", func(t *testing.T) {
		t.Run("writes a forbidden status", func(t *testing.T) {
			status, _ := callRequireRole("admin", "viewer")
			if status != http.StatusForbidden {
				t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
			}
		})

		// NOTE: this documents the current (unrefactored) behavior: the
		// forbidden response gets written, but the downstream handler still
		// runs right after it. Once RequireRole actually stops the chain on
		// a mismatch, this assertion should flip to false.
		t.Run("downstream handler currently still runs", func(t *testing.T) {
			_, downstream := callRequireRole("admin", "viewer")
			if !downstream.called {
				t.Fatal("expected downstream handler to currently still run despite the forbidden response (pre-refactor behavior)")
			}
		})
	})
}

func TestNewPayrollHandler(t *testing.T) {
	t.Run("admin role", func(t *testing.T) {
		t.Run("known employee", func(t *testing.T) {
			t.Run("returns an ok status", func(t *testing.T) {
				status, _ := callPayrollHandler("e-100", "admin")
				if status != http.StatusOK {
					t.Fatalf("status = %d, want %d", status, http.StatusOK)
				}
			})

			t.Run("adjusts the salary", func(t *testing.T) {
				_, body := callPayrollHandler("e-100", "admin")
				if !strings.Contains(body, "ok: e-100") {
					t.Fatalf("expected body to confirm the adjustment, got %q", body)
				}
			})
		})

		t.Run("unknown employee", func(t *testing.T) {
			t.Run("returns a not found status", func(t *testing.T) {
				status, _ := callPayrollHandler("e-999", "admin")
				if status != http.StatusNotFound {
					t.Fatalf("status = %d, want %d", status, http.StatusNotFound)
				}
			})
		})
	})

	t.Run("non-admin role", func(t *testing.T) {
		t.Run("writes a forbidden status", func(t *testing.T) {
			status, _ := callPayrollHandler("e-100", "viewer")
			if status != http.StatusForbidden {
				t.Fatalf("status = %d, want %d", status, http.StatusForbidden)
			}
		})

		// NOTE: same pre-refactor behavior as above, seen end to end -- a
		// non-admin caller still gets the salary adjusted even though the
		// response says forbidden.
		t.Run("salary adjustment currently still happens", func(t *testing.T) {
			_, body := callPayrollHandler("e-100", "viewer")
			if !strings.Contains(body, "ok: e-100") {
				t.Fatalf("expected the adjustment to currently still run despite the forbidden response (pre-refactor behavior), got %q", body)
			}
		})
	})
}
