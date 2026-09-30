package reminder

import (
	"strings"
	"testing"
	"time"
)

// These helpers are expected to change as you refactor — update them to match your new API.
func newTestService(host, from string) *ReminderService {
	return NewReminderService(host, 587, from)
}

func sendReminder(s *ReminderService, appt Appointment) error {
	return s.SendReminder(appt)
}

func sendBatch(s *ReminderService, appts []Appointment) []error {
	return s.SendBatch(appts)
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func testAppointment(email string) Appointment {
	return Appointment{
		Patient: Patient{ID: "p1", Name: "Jamie Rivera", Email: email},
		Time:    time.Date(2026, 10, 5, 9, 30, 0, 0, time.UTC),
		Doctor:  "Okafor",
		Clinic:  "Riverside Family Clinic",
	}
}

func TestSendReminder(t *testing.T) {
	t.Run("valid appointment with configured relay sends cleanly", func(t *testing.T) {
		svc := newTestService("smtp.clinicrelay.test", "reminders@riverside.test")
		err := sendReminder(svc, testAppointment("jamie@example.com"))
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	})

	t.Run("missing relay host returns an error", func(t *testing.T) {
		svc := newTestService("", "reminders@riverside.test")
		err := sendReminder(svc, testAppointment("jamie@example.com"))
		if err == nil {
			t.Fatal("expected an error when no relay host is configured")
		}
	})

	t.Run("missing relay host error references the patient", func(t *testing.T) {
		svc := newTestService("", "reminders@riverside.test")
		err := sendReminder(svc, testAppointment("jamie@example.com"))
		if err == nil || !strings.Contains(err.Error(), "p1") {
			t.Fatalf("expected error to reference patient id %q, got %v", "p1", err)
		}
	})

	t.Run("invalid patient email returns an error", func(t *testing.T) {
		svc := newTestService("smtp.clinicrelay.test", "reminders@riverside.test")
		err := sendReminder(svc, testAppointment("not-an-email"))
		if err == nil {
			t.Fatal("expected an error for an email address without an @")
		}
	})
}

func TestSendBatch(t *testing.T) {
	t.Run("all appointments succeed", func(t *testing.T) {
		svc := newTestService("smtp.clinicrelay.test", "reminders@riverside.test")
		appts := []Appointment{testAppointment("a@example.com"), testAppointment("b@example.com")}
		errs := sendBatch(svc, appts)
		if len(errs) != 0 {
			t.Fatalf("expected no errors, got %v", errs)
		}
	})

	t.Run("one appointment fails while the other succeeds", func(t *testing.T) {
		svc := newTestService("smtp.clinicrelay.test", "reminders@riverside.test")
		appts := []Appointment{testAppointment("a@example.com"), testAppointment("not-an-email")}
		errs := sendBatch(svc, appts)
		if len(errs) != 1 {
			t.Fatalf("expected exactly 1 error, got %d: %v", len(errs), errs)
		}
	})
}
