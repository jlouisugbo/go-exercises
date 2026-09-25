package snapshot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrInvalidCustomer      = errors.New("invalid customer")
	ErrCustomerNotFound     = errors.New("customer not found")
	ErrDirectoryUnavailable = errors.New("customer directory unavailable")
)

type Customer struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Order struct {
	ID         string `json:"id"`
	TotalCents int    `json:"total_cents"`
}

type Snapshot struct {
	Customer     Customer `json:"customer"`
	RecentOrders []Order  `json:"recent_orders"`
}

type Directory interface {
	Customer(ctx context.Context, customerID string) (Customer, error)
	RecentOrders(ctx context.Context, customerID string) ([]Order, error)
}

type Service struct {
	directory Directory
}

func NewService(directory Directory) *Service {
	return &Service{directory: directory}
}

func (s *Service) Build(customerID string) (Snapshot, error) {
	customerID = strings.TrimSpace(customerID)
	if customerID == "" {
		return Snapshot{}, ErrInvalidCustomer
	}

	customer, err := s.directory.Customer(context.Background(), customerID)
	if err != nil {
		return Snapshot{}, fmt.Errorf("load customer: %w", err)
	}

	orders, err := s.directory.RecentOrders(context.Background(), customerID)
	if err != nil {
		return Snapshot{}, fmt.Errorf("load recent orders: %w", err)
	}

	return Snapshot{Customer: customer, RecentOrders: orders}, nil
}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	snapshot, err := h.service.Build(r.URL.Query().Get("customer_id"))
	if err != nil {
		writeSnapshotError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(snapshot)
}

func writeSnapshotError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidCustomer):
		http.Error(w, ErrInvalidCustomer.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrCustomerNotFound):
		http.Error(w, ErrCustomerNotFound.Error(), http.StatusNotFound)
	case errors.Is(err, ErrDirectoryUnavailable):
		http.Error(w, ErrDirectoryUnavailable.Error(), http.StatusBadGateway)
	default:
		http.Error(w, "could not build customer snapshot", http.StatusInternalServerError)
	}
}
