package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gotest.tools/v3/assert"
)

func TestProvider(t *testing.T) {
	err := New().InternalValidate()
	assert.NilError(t, err)
}

func testAccProviders(t *testing.T) map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"infra": func() (*schema.Provider, error) {
			return New(), nil
		},
	}
}

func testAccPreCheck(t *testing.T) func() {
	return func() {
		accessKey := os.Getenv("INFRA_ACCESS_KEY")
		assert.Assert(t, accessKey != "", "`INFRA_ACCESS_KEY` must be set for acceptance tests")
	}
}
