package hashtags

import (
	"errors"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func Insert(tags []string, text string) string {
	return insertIntoHTML(tags, text)
}

func insertIntoHTML(tags []string, source string) string {
	z := html.NewTokenizer(strings.NewReader(source))
	hashtagAllowed := true
	var out string
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if errors.Is(z.Err(), io.EOF) {
				break
			} else {
				panic(z.Err())
			}
		} else if tt == html.TextToken && hashtagAllowed {
			text, newTags := insertIntoText(tags, string(z.Text()))
			out += text
			tags = newTags
		} else {
			if tt == html.StartTagToken {
				name, _ := z.TagName()
				if string(name) == "a" {
					hashtagAllowed = false
				}
			} else if tt == html.EndTagToken {
				name, _ := z.TagName()
				if string(name) == "a" {
					hashtagAllowed = true
				}
			}
			out += string(z.Raw())
		}
	}

	if len(tags) > 0 {
		var tagsWithHash []string
		for _, tag := range tags {
			tagsWithHash = append(tagsWithHash, "#"+tag)
		}
		out = strings.Join(tagsWithHash, " ") + "\n\n" + out
	}

	return out
}

func insertIntoText(tags []string, text string) (string, []string) {
	var missingTags []string
	for _, tag := range tags {
		tagRegexp := regexp.MustCompile(`\b` + tag + `\b`)
		match := tagRegexp.FindStringIndex(text)
		if match != nil {
			text = text[0:match[0]] + "#" + text[match[0]:]
		} else {
			missingTags = append(missingTags, tag)
		}
	}
	return text, missingTags
}
