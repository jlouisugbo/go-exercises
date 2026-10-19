package submissionfeed

import (
	"errors"
	"testing"
)

// newStore and submit are expected to change as you refactor — update them
// to match your new API.
func newStore() *Store {
	return NewStore()
}

func submit(store *Store, listeners []Listener, studentID, assignmentID string) error {
	return SubmitAssignment(store, listeners, studentID, assignmentID)
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestSubmitAssignment(t *testing.T) {
	t.Run("no listeners registered", func(t *testing.T) {
		store := newStore()

		err := submit(store, nil, "stu-1", "asn-1")

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("still saves the submission", func(t *testing.T) {
			if store.Count() != 1 {
				t.Errorf("expected 1 saved submission, got %d", store.Count())
			}
		})
	})

	t.Run("every listener succeeds", func(t *testing.T) {
		store := newStore()
		var calls []string
		listeners := []Listener{
			func(sub Submission) error {
				calls = append(calls, "gradebook")
				return nil
			},
			func(sub Submission) error {
				calls = append(calls, "email")
				return nil
			},
		}

		err := submit(store, listeners, "stu-2", "asn-1")

		t.Run("returns no error", func(t *testing.T) {
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
		t.Run("saves the submission", func(t *testing.T) {
			if store.Count() != 1 {
				t.Errorf("expected 1 saved submission, got %d", store.Count())
			}
		})
		t.Run("calls every listener in registration order", func(t *testing.T) {
			want := []string{"gradebook", "email"}
			if len(calls) != len(want) {
				t.Fatalf("expected %v, got %v", want, calls)
			}
			for i := range want {
				if calls[i] != want[i] {
					t.Errorf("expected %v, got %v", want, calls)
					break
				}
			}
		})
	})

	t.Run("a listener returns an error", func(t *testing.T) {
		store := newStore()
		var calls []string
		listeners := []Listener{
			func(sub Submission) error {
				calls = append(calls, "plagiarism-scan")
				return errors.New("scan service unavailable")
			},
			func(sub Submission) error {
				calls = append(calls, "gradebook")
				return nil
			},
		}

		err := submit(store, listeners, "stu-3", "asn-1")

		t.Run("returns an error", func(t *testing.T) {
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
		})
		t.Run("the submission is saved anyway", func(t *testing.T) {
			if store.Count() != 1 {
				t.Errorf("expected 1 saved submission, got %d", store.Count())
			}
		})
		// This pins today's behavior: a failing listener stops every
		// listener registered after it from running at all. That's
		// expected to change once listeners stop being able to block
		// each other.
		t.Run("the next listener never runs", func(t *testing.T) {
			if len(calls) != 1 {
				t.Errorf("expected only the failing listener to have run, got calls: %v", calls)
			}
		})
	})
}
