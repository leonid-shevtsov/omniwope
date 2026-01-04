package discord

import (
	"testing"

	"github.com/leonid-shevtsov/omniwope/internal/output/discord/api"
	"github.com/stretchr/testify/assert"
)

func TestCalculateChecksum(t *testing.T) {
	tests := []struct {
		name          string
		embed         api.Embed
		resourceBytes []byte
		wantHasColon  bool
	}{
		{
			name: "embed without resource",
			embed: api.Embed{
				Title:       "Test Title",
				Description: "Test Description",
			},
			resourceBytes: nil,
			wantHasColon:  false,
		},
		{
			name: "embed with resource",
			embed: api.Embed{
				Title:       "Test Title",
				Description: "Test Description",
			},
			resourceBytes: []byte("resource data"),
			wantHasColon:  true,
		},
		{
			name: "empty embed",
			embed: api.Embed{
				Title:       "",
				Description: "",
			},
			resourceBytes: nil,
			wantHasColon:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateChecksum(tt.embed, tt.resourceBytes)
			assert.NotEmpty(t, got)
			if tt.wantHasColon {
				assert.Contains(t, got, ":")
			} else {
				assert.NotContains(t, got, ":")
			}
		})
	}
}

func TestCalculateChecksumConsistency(t *testing.T) {
	embed := api.Embed{
		Title:       "Test Title",
		Description: "Test Description",
	}
	resourceBytes := []byte("test resource")

	// Same input should produce same checksum
	checksum1 := calculateChecksum(embed, resourceBytes)
	checksum2 := calculateChecksum(embed, resourceBytes)
	assert.Equal(t, checksum1, checksum2)

	// Different input should produce different checksum
	embed2 := api.Embed{
		Title:       "Different Title",
		Description: "Test Description",
	}
	checksum3 := calculateChecksum(embed2, resourceBytes)
	assert.NotEqual(t, checksum1, checksum3)
}
