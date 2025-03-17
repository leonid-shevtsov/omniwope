package linkparser

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
)

func TestParser(t *testing.T) {
	refTransformer := NewRefTransformer(
		func(doc string) string {
			return "https://ref/" + doc
		},
		func(doc string) string {
			return "https://relref/" + doc
		},
	)
	md := goldmark.New(goldmark.WithParserOptions(parser.WithASTTransformers(util.Prioritized(refTransformer, 0))))
	var out bytes.Buffer
	source := []byte("[foo]({{< relref \"foo.md\" >}}) [bar]({{< ref \"bar.md\" >}}) {{< ref \"bar\" >}}")
	preprocessedSource := PreprocessRefs(source)
	// nodes := md.Parser().Parse(text.NewReader(preprocessedSource))
	// nodes.Dump(source, 2)

	err := md.Convert(preprocessedSource, &out)

	log.Print(string(out.Bytes()))

	result := UndoRefs(out.Bytes())

	require.NoError(t, err)
	assert.Equal(t, `<p><a href="https://relref/foo.md">foo</a> <a href="https://ref/bar.md">bar</a> {{< ref "bar" >}}</p>`+"\n", string(result))
}
