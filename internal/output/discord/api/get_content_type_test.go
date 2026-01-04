package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContentType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "jpeg file",
			filename: "image.jpg",
			want:     "image/jpeg",
		},
		{
			name:     "jpeg uppercase",
			filename: "image.JPEG",
			want:     "image/jpeg",
		},
		{
			name:     "png file",
			filename: "image.png",
			want:     "image/png",
		},
		{
			name:     "gif file",
			filename: "animation.gif",
			want:     "image/gif",
		},
		{
			name:     "webp file",
			filename: "image.webp",
			want:     "image/webp",
		},
		{
			name:     "mp4 video",
			filename: "video.mp4",
			want:     "video/mp4",
		},
		{
			name:     "webm video",
			filename: "video.webm",
			want:     "video/webm",
		},
		{
			name:     "mov video",
			filename: "video.mov",
			want:     "video/quicktime",
		},
		{
			name:     "unknown extension",
			filename: "file.xyz",
			want:     "application/octet-stream",
		},
		{
			name:     "no extension",
			filename: "file",
			want:     "application/octet-stream",
		},
		{
			name:     "path with directory",
			filename: "/path/to/image.jpg",
			want:     "image/jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetContentType(tt.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}
