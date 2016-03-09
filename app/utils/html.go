package utils

import (
	"github.com/russross/blackfriday"
	"html/template"
	"regexp"
	"strings"
)

func Html2str(html string) string {
	src := string(html)

	// Lowercase HTML tags
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	// Remove styles
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	// Remove scripts
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	// Strip all HTML tags
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	// Remove continuous `\n`
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

func Markdown2Html(text string) string {
	return string(blackfriday.MarkdownCommon([]byte(text)))
}

func Markdown2HtmlTemplate(text string) template.HTML {
	return template.HTML(Markdown2Html(text))
}
