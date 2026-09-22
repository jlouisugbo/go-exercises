package workerpool

import (
	"reflect"
	"testing"
)

func TestProcessJobs(t *testing.T) {
	tests := []struct {
		name    string
		jobs    []int
		workers int
		want    []int
		wantErr bool
	}{
		{
			name:    "processes jobs in order",
			jobs:    []int{1, 2, 3, 4, 5},
			workers: 3,
			want:    []int{1, 4, 9, 16, 25},
		},
		{
			name:    "single worker",
			jobs:    []int{2, 4, 6},
			workers: 1,
			want:    []int{4, 16, 36},
		},
		{
			name:    "more workers than jobs",
			jobs:    []int{3, 5},
			workers: 8,
			want:    []int{9, 25},
		},
		{
			name:    "empty input",
			jobs:    []int{},
			workers: 2,
			want:    []int{},
		},
		{
			name:    "invalid worker count",
			jobs:    []int{1, 2, 3},
			workers: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProcessJobs(tt.jobs, tt.workers)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
