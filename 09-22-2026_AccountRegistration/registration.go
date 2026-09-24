package registration

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type User struct {
	Email string `json:"email"`
	Plan  Plan   `json:"plan"`
}

type Email struct {
	To      string
	Subject string
	Body    string
}

type Plan string

const (
	FreePlan Plan = "free"
	ProPlan  Plan = "pro"
	TeamPlan Plan = "team"
)

var (
	ErrUnavailable       = errors.New("dependency unavailable")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidPlan       = errors.New("invalid plan")
	ErrAlreadyRegistered = errors.New("email already registered")
	ErrSaveUser          = errors.New("could not save user")
	ErrSendEmail         = errors.New("could not send welcome email")
)

type UserStore interface {
	Has(email string) bool
	Save(user User) error
}

type EmailSender interface {
	Send(email Email) error
}

type MemoryUserStore struct {
	mu      sync.Mutex
	users   map[string]User
	SaveErr error
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{users: make(map[string]User)}
}

func (s *MemoryUserStore) Has(email string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.users[email]
	return ok
}

func (s *MemoryUserStore) Save(user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.SaveErr != nil {
		return s.SaveErr
	}
	s.users[user.Email] = user
	return nil
}

func (s *MemoryUserStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.users)
}

type MemoryEmailSender struct {
	mu      sync.Mutex
	Emails  []Email
	SendErr error
}

func (s *MemoryEmailSender) Send(email Email) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.SendErr != nil {
		return s.SendErr
	}
	s.Emails = append(s.Emails, email)
	return nil
}

func (s *MemoryEmailSender) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.Emails)
}

// Service is the one registration path. HTTP, a CSV import, or a queue
// consumer all call Register with an email and a plan.
type Service struct {
	users  UserStore
	emails EmailSender
}

func NewService(users UserStore, emails EmailSender) *Service {
	return &Service{users: users, emails: emails}
}

func (s *Service) Register(email string, plan Plan) (User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || !strings.Contains(email, "@") {
		return User{}, ErrInvalidEmail
	}
	if plan != FreePlan && plan != ProPlan && plan != TeamPlan {
		return User{}, ErrInvalidPlan
	}
	if s.users.Has(email) {
		return User{}, ErrAlreadyRegistered
	}

	user := User{Email: email, Plan: plan}
	if err := s.users.Save(user); err != nil {
		return User{}, fmt.Errorf("%w: %w", ErrSaveUser, err)
	}

	message := Email{
		To:      user.Email,
		Subject: "Welcome",
		Body:    fmt.Sprintf("Your %s account is ready.", user.Plan),
	}
	if err := s.emails.Send(message); err != nil {
		return User{}, fmt.Errorf("%w: %w", ErrSendEmail, err)
	}
	return user, nil
}

type RegistrationHandler struct {
	svc *Service
}

func NewRegistrationHandler(users UserStore, emails EmailSender) *RegistrationHandler {
	return &RegistrationHandler{svc: NewService(users, emails)}
}

func (h *RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Email string `json:"email"`
		Plan  Plan   `json:"plan"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.svc.Register(input.Email, input.Plan)
	if err != nil {
		writeRegistrationError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

func writeRegistrationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidEmail):
		http.Error(w, ErrInvalidEmail.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrInvalidPlan):
		http.Error(w, ErrInvalidPlan.Error(), http.StatusBadRequest)
	case errors.Is(err, ErrAlreadyRegistered):
		http.Error(w, ErrAlreadyRegistered.Error(), http.StatusConflict)
	case errors.Is(err, ErrSaveUser):
		http.Error(w, ErrSaveUser.Error(), http.StatusInternalServerError)
	case errors.Is(err, ErrSendEmail):
		http.Error(w, ErrSendEmail.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, "registration failed", http.StatusInternalServerError)
	}
}
