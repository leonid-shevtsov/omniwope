package discordmarkdown

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
			name: "bold",
			in:   "**bold**",
			out:  "**bold**\n\n",
		},
		{
			name: "italic",
			in:   "*italic*",
			out:  "*italic*\n\n",
		},
		{
			name: "bold italic",
			in:   "***bold italic***",
			out:  "***bold italic***\n\n",
		},
		{
			name: "code span",
			in:   "`code`",
			out:  "`code`\n\n",
		},
		{
			name: "fenced code block",
			in:   "```go\nfunc main() {}\n```",
			out:  "```go\nfunc main() {}\n```\n",
		},
		{
			name: "code block",
			in:   "    code block",
			out:  "```\ncode block\n```\n",
		},
		{
			name: "link",
			in:   "[text](https://example.com)",
			out:  "[text](https://example.com)\n\n",
		},
		{
			name: "auto link",
			in:   "https://example.com",
			out:  "https://example.com\n\n",
		},
		{
			name: "blockquote",
			in:   "> quote",
			out:  "> quote\n\n",
		},
		{
			name: "heading converted to bold",
			in:   "# Heading",
			out:  "**Heading**\n",
		},
		{
			name: "list",
			in:   "- item 1\n- item 2",
			out:  "- item 1\n- item 2\n",
		},
		{
			name: "thematic break",
			in:   "---",
			out:  "---\n",
		},
		{
			name: "paragraph",
			in:   "paragraph text",
			out:  "paragraph text\n\n",
		},
		{
			name: "mixed content",
			in:   "# Title\n\nThis is **bold** and *italic* text with a [link](https://example.com).",
			out:  "**Title**\nThis is **bold** and *italic* text with a [link](https://example.com).\n\n",
		},
	}

	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			var writer bytes.Buffer
			err := md.Convert([]byte(example.in), &writer)
			assert.NoError(t, err)
			out := writer.String()
			assert.Equal(t, example.out, out, "input: %q", example.in)
		})
	}
}
