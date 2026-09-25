package snapshot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubDirectory struct {
	customerFn    func(context.Context, string) (Customer, error)
	ordersFn      func(context.Context, string) ([]Order, error)
	customerCalls int
	orderCalls    int
}

func (s *stubDirectory) Customer(ctx context.Context, customerID string) (Customer, error) {
	s.customerCalls++
	return s.customerFn(ctx, customerID)
}

func (s *stubDirectory) RecentOrders(ctx context.Context, customerID string) ([]Order, error) {
	s.orderCalls++
	return s.ordersFn(ctx, customerID)
}

type testApp struct {
	handler   *Handler
	directory *stubDirectory
}

type snapshotResult struct {
	status        int
	snapshot      Snapshot
	customerCalls int
	orderCalls    int
}

// This is the construction and public-call seam. If your refactor changes
// constructors or method signatures, adapt these helpers before changing tests.
func newTestApp() *testApp {
	directory := &stubDirectory{
		customerFn: func(_ context.Context, customerID string) (Customer, error) {
			return Customer{
				ID:    customerID,
				Name:  "Joel",
				Email: "joel@example.com",
			}, nil
		},
		ordersFn: func(_ context.Context, _ string) ([]Order, error) {
			return []Order{{ID: "order-1", TotalCents: 4200}}, nil
		},
	}

	service := NewService(directory)
	return &testApp{
		handler:   NewHandler(service),
		directory: directory,
	}
}

func (a *testApp) request(ctx context.Context, method, target string) snapshotResult {
	request := httptest.NewRequest(method, target, nil).WithContext(ctx)
	response := httptest.NewRecorder()
	a.handler.ServeHTTP(response, request)

	var snapshot Snapshot
	_ = json.Unmarshal(response.Body.Bytes(), &snapshot)

	return snapshotResult{
		status:        response.Code,
		snapshot:      snapshot,
		customerCalls: a.directory.customerCalls,
		orderCalls:    a.directory.orderCalls,
	}
}

func TestSnapshotHandler(t *testing.T) {
	t.Run("request handling", func(t *testing.T) {
		t.Run("rejects methods other than GET", func(t *testing.T) {
			app := newTestApp()
			result := app.request(context.Background(), http.MethodPost, "/snapshot?customer_id=customer-1")
			if result.status != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want %d", result.status, http.StatusMethodNotAllowed)
			}
		})

		t.Run("rejects a blank customer id", func(t *testing.T) {
			app := newTestApp()
			result := app.request(context.Background(), http.MethodGet, "/snapshot?customer_id=%20%20")
			if result.status != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d", result.status, http.StatusBadRequest)
			}
		})
	})

	t.Run("successful snapshot", func(t *testing.T) {
		t.Run("returns the customer", func(t *testing.T) {
			app := newTestApp()
			result := app.request(context.Background(), http.MethodGet, "/snapshot?customer_id=customer-1")
			if result.status != http.StatusOK {
				t.Fatalf("status = %d, want %d", result.status, http.StatusOK)
			}
			if result.snapshot.Customer.ID != "customer-1" {
				t.Fatalf("customer id = %q, want %q", result.snapshot.Customer.ID, "customer-1")
			}
		})

		t.Run("returns recent orders", func(t *testing.T) {
			app := newTestApp()
			result := app.request(context.Background(), http.MethodGet, "/snapshot?customer_id=customer-1")
			if len(result.snapshot.RecentOrders) != 1 {
				t.Fatalf("order count = %d, want 1", len(result.snapshot.RecentOrders))
			}
		})
	})

	t.Run("dependency failures", func(t *testing.T) {
		t.Run("returns not found for an unknown customer", func(t *testing.T) {
			app := newTestApp()
			app.directory.customerFn = func(context.Context, string) (Customer, error) {
				return Customer{}, ErrCustomerNotFound
			}

			result := app.request(context.Background(), http.MethodGet, "/snapshot?customer_id=missing")
			if result.status != http.StatusNotFound {
				t.Fatalf("status = %d, want %d", result.status, http.StatusNotFound)
			}
		})

		t.Run("does not load orders after the customer lookup fails", func(t *testing.T) {
			app := newTestApp()
			app.directory.customerFn = func(context.Context, string) (Customer, error) {
				return Customer{}, ErrDirectoryUnavailable
			}

			result := app.request(context.Background(), http.MethodGet, "/snapshot?customer_id=customer-1")
			if result.orderCalls != 0 {
				t.Fatalf("order calls = %d, want 0", result.orderCalls)
			}
		})

		t.Run("returns bad gateway when orders are unavailable", func(t *testing.T) {
			app := newTestApp()
			app.directory.ordersFn = func(context.Context, string) ([]Order, error) {
				return nil, ErrDirectoryUnavailable
			}

			result := app.request(context.Background(), http.MethodGet, "/snapshot?customer_id=customer-1")
			if result.status != http.StatusBadGateway {
				t.Fatalf("status = %d, want %d", result.status, http.StatusBadGateway)
			}
		})
	})
}
