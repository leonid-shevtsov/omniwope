package telegold

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuin/goldmark"
)

func TestRenderer(t *testing.T) {
	md := goldmark.New(goldmark.WithRenderer(NewRenderer()))

	examples := []struct {
		name string
		in   string
		out  string
	}{
		{
			"italic",
			"**italic**",
			"<b>italic</b>\n\n",
		},
		{
			"bold",
			"**bold**",
			"<b>bold</b>\n\n",
		},
	}

	for _, example := range examples {
		var writer bytes.Buffer
		md.Convert([]byte(example.in), &writer)
		out := writer.String()
		assert.Equal(t, example.out, out, example.name)
	}
}

// TODO:
// test case for tele gold - empty line between list and code block (or paragraph)
