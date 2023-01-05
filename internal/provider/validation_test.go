package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"gotest.tools/v3/assert"
)

func TestStringIsDuration(t *testing.T) {
	cases := map[string]diag.Diagnostics{
		"1h":      nil,
		"1h30m2s": nil,
		"30": diag.Diagnostics{
			{Summary: `time: missing unit in duration "30"`},
		},
		"abc": diag.Diagnostics{
			{Summary: `time: invalid duration "abc"`},
		},
	}

	fn := validateStringIsDuration()
	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			actual := fn(input, cty.Path{})
			assert.DeepEqual(t, actual, expected)
		})
	}
}

func TestStringIsEmail(t *testing.T) {
	cases := map[string]diag.Diagnostics{
		"admin@example.com": nil,
		"admin@example":     nil,
		"admin": diag.Diagnostics{
			{Summary: "mail: missing '@' or angle-addr"},
		},
		"admin <admin@example.com>": diag.Diagnostics{
			{Summary: "mail: must not contain display name"},
		},
	}

	fn := validateStringIsEmail()
	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			actual := fn(input, cty.Path{})
			assert.DeepEqual(t, actual, expected)
		})
	}
}

func TestStringIsID(t *testing.T) {
	cases := map[string]diag.Diagnostics{
		"1":          nil,
		"if":         nil,
		"gjLJvMroqy": nil,
		"0": diag.Diagnostics{
			{Summary: "invalid base58: byte 0 is out of range"},
		},
		"gjLJvMroqygjLJvMroqy": diag.Diagnostics{
			{Summary: "invalid base58: too long"},
		},
	}

	fn := validateStringIsID()
	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			actual := fn(input, cty.Path{})
			assert.DeepEqual(t, actual, expected)
		})
	}
}

func TestStringIsName(t *testing.T) {
	cases := map[string]diag.Diagnostics{
		"name":    nil,
		"nAmE":    nil,
		"n4m3":    nil,
		"na_-.me": nil,
		"": diag.Diagnostics{
			{Summary: "invalid value for  (Name may contain only letters (uppercase and lowercase), numbers, underscores `_`, hyphens `-`, and periods `.`.)", AttributePath: cty.Path{}},
			{Summary: "expected length of  to be in the range (2 - 256), got ", AttributePath: cty.Path{}},
		},
		"na@me": diag.Diagnostics{
			{Summary: "invalid value for  (Name may contain only letters (uppercase and lowercase), numbers, underscores `_`, hyphens `-`, and periods `.`.)", AttributePath: cty.Path{}},
		},
		"n": diag.Diagnostics{
			{Summary: "expected length of  to be in the range (2 - 256), got n", AttributePath: cty.Path{}},
		},
		strings.Repeat("0123456789abcdef", 17): diag.Diagnostics{
			{Summary: "expected length of  to be in the range (2 - 256), got " + strings.Repeat("0123456789abcdef", 17), AttributePath: cty.Path{}},
		},
	}

	fn := validateStringIsName()
	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			actual := fn(input, cty.Path{})
			assert.DeepEqual(t, actual, expected)
		})
	}
}

func TestStringMinLength(t *testing.T) {
	cases := map[string]diag.Diagnostics{
		"01234567":   nil,
		"0123456789": nil,
		"": diag.Diagnostics{
			{Summary: "expected length for  is at least 8", AttributePath: cty.Path{}},
		},
		"0": diag.Diagnostics{
			{Summary: "expected length for  is at least 8", AttributePath: cty.Path{}},
		},
		"0123456": diag.Diagnostics{
			{Summary: "expected length for  is at least 8", AttributePath: cty.Path{}},
		},
	}

	fn := StringMinLength(8)
	for input, expected := range cases {
		t.Run(input, func(t *testing.T) {
			actual := fn(input, cty.Path{})
			assert.DeepEqual(t, actual, expected)
		})
	}
}
