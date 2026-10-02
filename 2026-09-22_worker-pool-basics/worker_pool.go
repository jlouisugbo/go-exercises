package workerpool

import "errors"

// ProcessJobs processes jobs concurrently using a fixed number of workers.
//
// TODO:
//   1. Return an error if workers <= 0.
//   2. Start exactly workers goroutines.
//   3. Square each input value.
//   4. Preserve input ordering in the returned slice.
//   5. Avoid leaking goroutines.
func ProcessJobs(jobs []int, workers int) ([]int, error) {
	if workers <= 0 {
		return nil, errors.New("workers must be greater than zero")
	}

	// Your implementation here.
	return nil, nil
}
