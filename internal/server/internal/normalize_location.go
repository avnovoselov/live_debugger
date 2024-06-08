package internal

import (
	"fmt"
	"strings"
)

const (
	normalizeTrimLeftCutSet  = "/"
	normalizeLocationPattern = "/%s"
)

// NormalizeLocation - formats location to common format
//
//	/foo/baz/bar -> /foo/baz/bar
//	foo/baz/bar  -> /foo/bar
//	foo/baz/bar/ -> /foo/baz/bar
func NormalizeLocation(s string) string {
	return fmt.Sprintf(normalizeLocationPattern, strings.Trim(s, normalizeTrimLeftCutSet))
}
