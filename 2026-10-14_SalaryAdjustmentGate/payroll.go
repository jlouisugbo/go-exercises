// Package payroll exposes the internal HTTP endpoint payroll admins use
// to apply one-off salary adjustments.
package payroll

import (
	"fmt"
	"log"
	"net/http"
)

// Employee is a single row in the payroll system.
type Employee struct {
	ID          string
	Name        string
	SalaryCents int64
}

// store is a stand-in for the real payroll database.
type store struct {
	employees map[string]*Employee
}

func newStore() *store {
	return &store{
		employees: map[string]*Employee{
			"e-100": {ID: "e-100", Name: "Dana Whitfield", SalaryCents: 8200000},
			"e-101": {ID: "e-101", Name: "Marcus Oyelaran", SalaryCents: 6450000},
		},
	}
}

func (s *store) adjustSalary(id string, deltaCents int64) (*Employee, error) {
	emp, ok := s.employees[id]
	if !ok {
		return nil, fmt.Errorf("employee %s not found", id)
	}
	emp.SalaryCents += deltaCents
	return emp, nil
}

// AdjustSalaryHandler applies a flat one-off adjustment to the named
// employee. It's meant to be reached only through RequireRole, below --
// on its own it has no idea who's calling it.
func AdjustSalaryHandler(s *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("employee_id")

		emp, err := s.adjustSalary(id, 500000) // flat $5,000 bump for this endpoint
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		log.Printf("adjusted salary for %s to %d cents", emp.ID, emp.SalaryCents)
		fmt.Fprintf(w, "ok: %s new salary %d cents\n", emp.ID, emp.SalaryCents)
	}
}

// RequireRole only lets a request reach next if the caller's
// X-Employee-Role header matches role. Anything else is supposed to get
// rejected before it can touch payroll data.
func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Employee-Role") != role {
			log.Printf("rejected request from role %q, want %q", r.Header.Get("X-Employee-Role"), role)
			http.Error(w, "forbidden: admin role required", http.StatusForbidden)
		}

		next.ServeHTTP(w, r)
	})
}

// WithRequestLog logs every request that passes through the chain before
// it reaches the next handler.
func WithRequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// NewPayrollHandler wires up the payroll adjustment endpoint: every
// request is logged, then checked for the admin role, before it's allowed
// anywhere near AdjustSalaryHandler.
func NewPayrollHandler() http.Handler {
	s := newStore()
	return WithRequestLog(RequireRole("admin", AdjustSalaryHandler(s)))
}
