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
	Plan  string `json:"plan"`
}

type Email struct {
	To      string
	Subject string
	Body    string
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

type RegistrationHandler struct {
	Users  *MemoryUserStore
	Emails *MemoryEmailSender
}

func NewRegistrationHandler(users *MemoryUserStore, emails *MemoryEmailSender) *RegistrationHandler {
	return &RegistrationHandler{Users: users, Emails: emails}
}

func (h *RegistrationHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Email string `json:"email"`
		Plan  string `json:"plan"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" || !strings.Contains(email, "@") {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}

	if input.Plan != "free" && input.Plan != "pro" && input.Plan != "team" {
		http.Error(w, "invalid plan", http.StatusBadRequest)
		return
	}

	if h.Users.Has(email) {
		http.Error(w, "email already registered", http.StatusConflict)
		return
	}

	user := User{Email: email, Plan: input.Plan}
	if err := h.Users.Save(user); err != nil {
		http.Error(w, "could not save user", http.StatusInternalServerError)
		return
	}

	message := Email{
		To:      user.Email,
		Subject: "Welcome",
		Body:    fmt.Sprintf("Your %s account is ready.", user.Plan),
	}
	if err := h.Emails.Send(message); err != nil {
		http.Error(w, "could not send welcome email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

var ErrUnavailable = errors.New("dependency unavailable")
