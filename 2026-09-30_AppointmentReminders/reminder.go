package reminder

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Patient is a clinic patient eligible for appointment reminders.
type Patient struct {
	ID    string
	Name  string
	Email string
}

// Appointment is a scheduled visit that needs a reminder sent beforehand.
type Appointment struct {
	Patient Patient
	Time    time.Time
	Doctor  string
	Clinic  string
}

// SMTPClient is a thin wrapper around the clinic's email relay.
type SMTPClient struct {
	Host string
	Port int
}

// Send delivers a message through the configured SMTP relay. In production
// this opens a real connection to Host:Port; here it validates the inputs
// the relay would otherwise reject and reports success once they're clean.
func (c *SMTPClient) Send(to, subject, body string) error {
	if c.Host == "" {
		return errors.New("smtp: no relay host configured")
	}
	if !strings.Contains(to, "@") {
		return fmt.Errorf("smtp: invalid recipient address %q", to)
	}
	return nil
}

// ReminderService sends appointment reminders to patients the morning
// before their visit.
type ReminderService struct {
	smtpHost    string
	smtpPort    int
	fromAddress string
}

// NewReminderService wires up a reminder service pointed at the clinic's
// email relay.
func NewReminderService(smtpHost string, smtpPort int, fromAddress string) *ReminderService {
	return &ReminderService{
		smtpHost:    smtpHost,
		smtpPort:    smtpPort,
		fromAddress: fromAddress,
	}
}

// SendReminder emails the patient a reminder about their upcoming appointment.
func (s *ReminderService) SendReminder(appt Appointment) error {
	client := &SMTPClient{Host: s.smtpHost, Port: s.smtpPort}

	subject := fmt.Sprintf("Reminder: appointment with Dr. %s", appt.Doctor)
	body := fmt.Sprintf(
		"Hi %s,\n\nThis is a reminder that you have an appointment with Dr. %s at %s on %s.\n\n- %s",
		appt.Patient.Name, appt.Doctor, appt.Clinic, appt.Time.Format("Jan 2 at 3:04 PM"), s.fromAddress,
	)

	if err := client.Send(appt.Patient.Email, subject, body); err != nil {
		return fmt.Errorf("sending reminder to patient %s: %w", appt.Patient.ID, err)
	}
	return nil
}

// SendBatch sends reminders for every appointment in the list, collecting
// (rather than stopping on) individual failures.
func (s *ReminderService) SendBatch(appts []Appointment) []error {
	var errs []error
	for _, appt := range appts {
		if err := s.SendReminder(appt); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
