package internal

import (
	"fmt"
	"strings"
)

const (
	normalizeTrimLeftCutSet  = "/"
	normalizeLocationPattern = "/%s"
)

func NormalizeLocation(s string) string {
	return fmt.Sprintf(normalizeLocationPattern, strings.TrimLeft(s, normalizeTrimLeftCutSet))
}
