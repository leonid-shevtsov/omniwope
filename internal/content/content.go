package content

import (
	"encoding/json"
	"time"

	"github.com/leonid-shevtsov/omniwope/internal/checksum"
)

type Post struct {
	// URL is used as the unique identifier of the post
	URL       string     `json:"url"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Date      time.Time  `json:"date"`
	Resources []Resource `json:"resources"`
	Tags      []string   `json:"tags"`
}

const resourceTypeImage = "image"
const resourceTypeVideo = "video"

type Resource struct {
	Label string `json:"label"`
	Path  string `json:"path"`
	Type  string `json:"type"`
}

// For now just do a plain JSON encoding.
// Outputs should do a checksum of transformed content as well
func (p Post) Checksum() string {
	encoded, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	return checksum.Sum(encoded)
}

func (r Resource) IsImage() bool {
	return r.Type == resourceTypeImage
}

func (r Resource) IsVideo() bool {
	return r.Type == resourceTypeVideo
}
