package compile

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

func toCamel(snake string) string {
	c := cases.Title(language.Dutch)
	items := strings.Split(snake, "_")
	for i, _ := range items {
		items[i] = c.String(items[i])
	}
	return strings.Join(items, "")
}
