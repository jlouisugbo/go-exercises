package shipment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Shipment represents a package somewhere between a warehouse and a
// customer's door.
type Shipment struct {
	ID          string     `json:"id"`
	Status      string     `json:"status"`
	CarrierETA  time.Time  `json:"carrier_eta"`
	Escalated   bool       `json:"escalated"`
	EscalatedAt *time.Time `json:"escalated_at,omitempty"`
}

// Server holds the in-memory shipment store and serves the status-update
// endpoint.
type Server struct {
	shipments map[string]*Shipment
}

// NewServer creates a Server seeded with the given shipments.
func NewServer(seed map[string]*Shipment) *Server {
	return &Server{shipments: seed}
}

type updateRequest struct {
	ShipmentID string `json:"shipment_id"`
	NewStatus  string `json:"new_status"`
}

var validStatuses = map[string]bool{
	"pending":    true,
	"in_transit": true,
	"delayed":    true,
	"delivered":  true,
}

// escalationGracePeriod is how long a shipment can sit past its carrier ETA
// before ops wants to be paged about it.
const escalationGracePeriod = 2 * time.Hour

// UpdateStatusHandler handles POST /shipments/status. It decodes the
// request body, validates it, looks up the shipment, applies the status
// transition (including the escalation rule for shipments that go delayed
// well past their carrier ETA), and writes the updated shipment back as
// JSON.
func (s *Server) UpdateStatusHandler(w http.ResponseWriter, r *http.Request) {
	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	req.ShipmentID = strings.TrimSpace(req.ShipmentID)
	req.NewStatus = strings.TrimSpace(req.NewStatus)

	if req.ShipmentID == "" {
		http.Error(w, "shipment_id is required", http.StatusBadRequest)
		return
	}

	if !validStatuses[req.NewStatus] {
		http.Error(w, fmt.Sprintf("unknown status %q", req.NewStatus), http.StatusBadRequest)
		return
	}

	sh, ok := s.shipments[req.ShipmentID]
	if !ok {
		http.Error(w, "shipment not found", http.StatusNotFound)
		return
	}

	now := time.Now()
	if req.NewStatus == "delayed" {
		if now.Sub(sh.CarrierETA) > escalationGracePeriod {
			sh.Escalated = true
			escalatedAt := now
			sh.EscalatedAt = &escalatedAt
		}
	} else if req.NewStatus == "delivered" {
		sh.Escalated = false
		sh.EscalatedAt = nil
	}

	sh.Status = req.NewStatus

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(sh)
}
