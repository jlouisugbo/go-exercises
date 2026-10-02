package shipment

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// These helpers are expected to change as you refactor — update them to match your new API.
func newTestServer(seed map[string]*Shipment) *Server {
	return NewServer(seed)
}

func doUpdateStatus(t *testing.T, s *Server, shipmentID, newStatus string) (*httptest.ResponseRecorder, Shipment) {
	t.Helper()
	body, err := json.Marshal(updateRequest{ShipmentID: shipmentID, NewStatus: newStatus})
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/shipments/status", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	s.UpdateStatusHandler(rec, req)

	var got Shipment
	if rec.Code == http.StatusOK {
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}
	return rec, got
}

func doRawUpdateStatus(t *testing.T, s *Server, rawBody string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/shipments/status", strings.NewReader(rawBody))
	rec := httptest.NewRecorder()
	s.UpdateStatusHandler(rec, req)
	return rec
}

// End of helpers — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestUpdateStatusHandler_InvalidBody(t *testing.T) {
	s := newTestServer(map[string]*Shipment{})

	t.Run("malformed JSON returns 400", func(t *testing.T) {
		rec := doRawUpdateStatus(t, s, "{not-json")
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}

func TestUpdateStatusHandler_MissingShipmentID(t *testing.T) {
	s := newTestServer(map[string]*Shipment{})

	t.Run("blank shipment_id returns 400", func(t *testing.T) {
		rec, _ := doUpdateStatus(t, s, "   ", "delayed")
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}

func TestUpdateStatusHandler_UnknownStatus(t *testing.T) {
	seed := map[string]*Shipment{
		"ship-1": {ID: "ship-1", Status: "pending", CarrierETA: time.Now()},
	}
	s := newTestServer(seed)

	t.Run("unrecognized status returns 400", func(t *testing.T) {
		rec, _ := doUpdateStatus(t, s, "ship-1", "lost_in_space")
		if rec.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}

func TestUpdateStatusHandler_ShipmentNotFound(t *testing.T) {
	s := newTestServer(map[string]*Shipment{})

	t.Run("unknown shipment id returns 404", func(t *testing.T) {
		rec, _ := doUpdateStatus(t, s, "ghost", "delayed")
		if rec.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusNotFound)
		}
	})
}

func TestUpdateStatusHandler_DelayedWithinGracePeriod(t *testing.T) {
	seed := map[string]*Shipment{
		"ship-1": {ID: "ship-1", Status: "in_transit", CarrierETA: time.Now().Add(-30 * time.Minute)},
	}
	s := newTestServer(seed)
	rec, got := doUpdateStatus(t, s, "ship-1", "delayed")

	t.Run("returns 200", func(t *testing.T) {
		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})
	t.Run("status becomes delayed", func(t *testing.T) {
		if got.Status != "delayed" {
			t.Errorf("status = %q, want %q", got.Status, "delayed")
		}
	})
	t.Run("not escalated within grace period", func(t *testing.T) {
		if got.Escalated {
			t.Errorf("escalated = true, want false")
		}
	})
}

func TestUpdateStatusHandler_DelayedPastGracePeriod(t *testing.T) {
	seed := map[string]*Shipment{
		"ship-1": {ID: "ship-1", Status: "in_transit", CarrierETA: time.Now().Add(-3 * time.Hour)},
	}
	s := newTestServer(seed)
	_, got := doUpdateStatus(t, s, "ship-1", "delayed")

	t.Run("escalated past grace period", func(t *testing.T) {
		if !got.Escalated {
			t.Errorf("escalated = false, want true")
		}
	})
	t.Run("escalated_at is set", func(t *testing.T) {
		if got.EscalatedAt == nil {
			t.Errorf("escalated_at = nil, want non-nil")
		}
	})
}

func TestUpdateStatusHandler_DeliveredClearsEscalation(t *testing.T) {
	escalatedAt := time.Now().Add(-1 * time.Hour)
	seed := map[string]*Shipment{
		"ship-1": {
			ID:          "ship-1",
			Status:      "delayed",
			CarrierETA:  time.Now().Add(-3 * time.Hour),
			Escalated:   true,
			EscalatedAt: &escalatedAt,
		},
	}
	s := newTestServer(seed)
	_, got := doUpdateStatus(t, s, "ship-1", "delivered")

	t.Run("status becomes delivered", func(t *testing.T) {
		if got.Status != "delivered" {
			t.Errorf("status = %q, want %q", got.Status, "delivered")
		}
	})
	t.Run("escalation cleared", func(t *testing.T) {
		if got.Escalated {
			t.Errorf("escalated = true, want false")
		}
	})
	t.Run("escalated_at cleared", func(t *testing.T) {
		if got.EscalatedAt != nil {
			t.Errorf("escalated_at = %v, want nil", got.EscalatedAt)
		}
	})
}

func TestUpdateStatusHandler_OtherStatusLeavesEscalationUntouched(t *testing.T) {
	escalatedAt := time.Now().Add(-1 * time.Hour)
	seed := map[string]*Shipment{
		"ship-1": {
			ID:          "ship-1",
			Status:      "delayed",
			CarrierETA:  time.Now().Add(-3 * time.Hour),
			Escalated:   true,
			EscalatedAt: &escalatedAt,
		},
	}
	s := newTestServer(seed)
	_, got := doUpdateStatus(t, s, "ship-1", "in_transit")

	t.Run("status updates", func(t *testing.T) {
		if got.Status != "in_transit" {
			t.Errorf("status = %q, want %q", got.Status, "in_transit")
		}
	})
	t.Run("escalation flag untouched by unrelated transition", func(t *testing.T) {
		if !got.Escalated {
			t.Errorf("escalated = false, want true (should be untouched)")
		}
	})
}
