package submissionfeed

import (
	"fmt"
	"time"
)

// Submission represents a single assignment turned in by a student.
type Submission struct {
	StudentID    string
	AssignmentID string
	SubmittedAt  time.Time
}

// Listener is notified whenever a submission comes in. Plagiarism scanning,
// gradebook syncing, and the "submission received" email to the student are
// all wired in as listeners on this same list.
type Listener func(Submission) error

// Store is where submissions actually live. It's deliberately tiny --
// everything worth looking at here is in how a submission gets announced,
// not in how it's stored.
type Store struct {
	saved []Submission
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{}
}

// Save appends a submission to the store.
func (s *Store) Save(sub Submission) {
	s.saved = append(s.saved, sub)
}

// Count returns how many submissions have been saved.
func (s *Store) Count() int {
	return len(s.saved)
}

// SubmitAssignment records the submission, then lets every registered
// listener know it happened, in the order they were registered.
func SubmitAssignment(store *Store, listeners []Listener, studentID, assignmentID string) error {
	sub := Submission{
		StudentID:    studentID,
		AssignmentID: assignmentID,
		SubmittedAt:  time.Now(),
	}

	store.Save(sub)

	for _, notify := range listeners {
		if err := notify(sub); err != nil {
			return fmt.Errorf("notifying listeners for %s: %w", studentID, err)
		}
	}

	return nil
}
