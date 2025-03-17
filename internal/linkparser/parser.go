package linkparser

import (
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

type refTransformer struct {
	linkHelper    func(string) string
	relLinkHelper func(string) string
}

var refOpener = []byte("{{<")
var refRegexp = regexp.MustCompile(`{{<\s*((?:rel)?)ref\s+"([^"]+)"\s*>}}`)
var refReplacement = []byte("REFSHORTCODE/${1}/${2}/")
var refReplacementRegexp = regexp.MustCompile(`^REFSHORTCODE/(.*)/(.+)/$`)
var refUndoRegexp = regexp.MustCompile(`REFSHORTCODE/([^/]*)/([^/]+)/`)
var refUndoReplacement = []byte(`{{&lt; ${1}ref &quot;${2}&quot; &gt;}}`)

func NewRefTransformer(linkHelper func(string) string, relLinkHelper func(string) string) parser.ASTTransformer {
	return &refTransformer{linkHelper, relLinkHelper}
}

func (p *refTransformer) Transform(node *ast.Document, reader text.Reader, pc parser.Context) {
	ast.Walk(node, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if linkNode, ok := n.(*ast.Link); ok {
			if match := refReplacementRegexp.FindSubmatch(linkNode.Destination); match != nil {
				isRel := len(match[1]) > 0
				var href string
				if isRel {
					href = p.relLinkHelper(string(match[2]))
				} else {
					href = p.linkHelper(string(match[2]))
				}
				linkNode.Destination = []byte(href)
			}
		}
		return ast.WalkContinue, nil
	})
}

func PreprocessRefs(input []byte) []byte {
	return refRegexp.ReplaceAll(input, refReplacement)
}

func UndoRefs(input []byte) []byte {
	return refUndoRegexp.ReplaceAll(input, refUndoReplacement)
}
