// Package thumbnails generates preview images for a batch of freshly
// uploaded photos.
package thumbnails

import "sync"

// Photo is a single uploaded image waiting to be thumbnailed.
type Photo struct {
	ID  string
	URL string
}

// ThumbnailResult is what GenerateThumbnails produces for one photo in
// the batch.
type ThumbnailResult struct {
	PhotoID   string
	Thumbnail string
	Err       error
}

// Resizer talks to the image-resizing service that actually produces a
// thumbnail for a photo.
type Resizer interface {
	Resize(photo Photo) (thumbnailRef string, err error)
}

// GenerateThumbnails produces a thumbnail for every photo in the batch.
// The order of the results doesn't matter -- callers look results up by
// PhotoID.
func GenerateThumbnails(photos []Photo, resizer Resizer) []ThumbnailResult {
	if len(photos) == 0 {
		return nil
	}

	var (
		mu      sync.Mutex
		wg      sync.WaitGroup
		results []ThumbnailResult
	)

	for _, photo := range photos {
		wg.Add(1)
		go func(p Photo) {
			defer wg.Done()

			thumb, err := resizer.Resize(p)

			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				results = append(results, ThumbnailResult{PhotoID: p.ID, Err: err})
				return
			}
			results = append(results, ThumbnailResult{PhotoID: p.ID, Thumbnail: thumb})
		}(photo)
	}

	wg.Wait()
	return results
}
