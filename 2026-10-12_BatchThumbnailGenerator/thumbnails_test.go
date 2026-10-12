package thumbnails

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// cappedResizer stands in for the real image-resizing service, which can
// only handle a limited number of requests at the same time. A call that
// arrives while the service is already at capacity gets rejected, the
// same way the real service would reject it with a 429.
type cappedResizer struct {
	max      int32
	delay    time.Duration
	inFlight int32
}

func (c *cappedResizer) Resize(p Photo) (string, error) {
	n := atomic.AddInt32(&c.inFlight, 1)
	defer atomic.AddInt32(&c.inFlight, -1)

	time.Sleep(c.delay)

	if n > c.max {
		return "", errors.New("resizer: capacity exceeded")
	}
	return "thumb-" + p.ID, nil
}

// This helper is expected to change as you refactor — update it to match your new API.
func runBatch(photos []Photo, resizer Resizer) []ThumbnailResult {
	return GenerateThumbnails(photos, resizer)
}

// End of helper — the tests below should not need to change as you refactor. But you are welcome to change them if you find you need to!

func TestGenerateThumbnails(t *testing.T) {
	t.Run("empty batch", func(t *testing.T) {
		results := runBatch(nil, &cappedResizer{max: 5})

		if len(results) != 0 {
			t.Fatalf("expected no results for an empty batch, got %d", len(results))
		}
	})

	t.Run("batch within capacity succeeds", func(t *testing.T) {
		resizer := &cappedResizer{max: 5, delay: 5 * time.Millisecond}
		photos := []Photo{{ID: "p1", URL: "http://example.com/p1.jpg"}}

		results := runBatch(photos, resizer)

		if len(results) != 1 {
			t.Fatalf("expected 1 result, got %d", len(results))
		}
		if results[0].Err != nil {
			t.Fatalf("expected no error, got %v", results[0].Err)
		}
	})

	t.Run("successful result carries a thumbnail reference", func(t *testing.T) {
		resizer := &cappedResizer{max: 5, delay: 5 * time.Millisecond}
		photos := []Photo{{ID: "p1", URL: "http://example.com/p1.jpg"}}

		results := runBatch(photos, resizer)

		if results[0].Thumbnail == "" {
			t.Fatalf("expected a thumbnail reference, got an empty string")
		}
	})

	// NOTE: this test documents the current (unrefactored) behavior: every
	// photo in the batch gets its own goroutine immediately, with nothing
	// limiting how many hit the resizer at once. Once the batch is driven
	// through a bounded number of workers (at or under the resizer's
	// capacity), this assertion should flip -- a batch this size should
	// come back with zero errors instead of several.
	t.Run("batch exceeding capacity currently produces rejections", func(t *testing.T) {
		resizer := &cappedResizer{max: 2, delay: 20 * time.Millisecond}
		photos := make([]Photo, 10)
		for i := range photos {
			photos[i] = Photo{ID: string(rune('a' + i)), URL: "http://example.com/x.jpg"}
		}

		results := runBatch(photos, resizer)

		if len(results) != len(photos) {
			t.Fatalf("expected a result for every photo, got %d for %d photos", len(results), len(photos))
		}

		var rejected int
		for _, r := range results {
			if r.Err != nil {
				rejected++
			}
		}
		if rejected == 0 {
			t.Fatalf("expected the unbounded fan-out to exceed the resizer's capacity and produce at least one rejection, got none")
		}
	})
}
