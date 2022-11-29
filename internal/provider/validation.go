package provider

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/infrahq/infra/uid"
)

func validateStringIsID() schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		if _, err := uid.Parse([]byte(v.(string))); err != nil {
			return diag.FromErr(err)
		}

		var diags diag.Diagnostics
		return diags
	}
}

func validateStringIsEmail() schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		address, err := mail.ParseAddress(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		if address.Name != "" {
			return diag.Errorf("mail: must not contain display name")
		}

		if !strings.Contains(strings.Split(address.Address, "@")[1], ".") {
			return diag.Errorf("mail: missing '.' in address domain")
		}

		var diags diag.Diagnostics
		return diags
	}
}

func validateStringIsName() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(
		validation.All(
			validation.StringMatch(regexp.MustCompile("^[\\w-.]+$"), "Name may contain only letters (uppercase and lowercase), numbers, underscores `_`, hyphens `-`, and periods `.`."),
			validation.StringLenBetween(2, 256),
		),
	)
}

func validateStringIsDuration() schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		if _, err := time.ParseDuration(v.(string)); err != nil {
			return diag.FromErr(err)
		}

		var diags diag.Diagnostics
		return diags
	}
}

func StringMinLength(min int) schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(
		func(v any, k string) (warnings []string, errors []error) {
			if len(v.(string)) < min {
				errors = append(errors, fmt.Errorf("expected length for %s is at least %d", k, min))
			}

			return warnings, errors
		},
	)
}

func validateStringIsPEMEncoded() schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		if _, err := DecodePEM([]byte(v.(string)), "CERTIFICATE"); err != nil {
			return diag.Errorf("failed to decode PEM block containing public key")
		}

		var diags diag.Diagnostics
		return diags
	}
}

func validateStringIsPEMEncodedFile() schema.SchemaValidateDiagFunc {
	return func(v any, p cty.Path) diag.Diagnostics {
		if _, err := DecodePEMFile(v.(string), "CERTIFICATE"); err != nil {
			return diag.FromErr(err)
		}

		var diags diag.Diagnostics
		return diags
	}
}
