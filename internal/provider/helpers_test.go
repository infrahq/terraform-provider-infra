package provider

import (
	"strings"
)

func composeTestConfigFunc(configs ...string) string {
	return strings.Join(configs, "\n")
}
