package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/infrahq/infra/api"
)

func resourceSettings() *schema.Resource {
	return &schema.Resource{
		Description: `Provides Infra organization settings.

` + "`infra_settings`" + ` behaves differently than normal Terraform resources as settings are
created with the organization. When a Terraform resource is created, settings automatically
imported while no action is taken when the resource is deleted.`,

		CreateContext: resourceSettingsUpdate,
		ReadContext:   resourceSettingsRead,
		UpdateContext: resourceSettingsUpdate,
		DeleteContext: resourceSettingsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"password_requirements": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum_length": {
							Description: "Minimum password length.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     8,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.IntAtLeast(8),
							),
						},
						"minimum_lowercase": {
							Description: "Minimum number of lowercase ASCII letters.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
						},
						"minimum_uppercase": {
							Description: "Minimum number of uppercase ASCII letters.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
						},
						"minimum_numbers": {
							Description: "Minimum number of numbers.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
						},
						"minimum_symbols": {
							Description: "Minimum number of symbols.",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
						},
					},
				},
			},
		},
	}
}

func resourceSettingsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	settings, err := client.GetSettings(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	requirements := map[string]int{
		"minimum_length":    settings.PasswordRequirements.LengthMin,
		"minimum_lowercase": settings.PasswordRequirements.LowercaseMin,
		"minimum_uppercase": settings.PasswordRequirements.UppercaseMin,
		"minimum_numbers":   settings.PasswordRequirements.NumberMin,
		"minimum_symbols":   settings.PasswordRequirements.SymbolMin,
	}

	if err := d.Set("password_requirements", []map[string]int{requirements}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("self")

	var diags diag.Diagnostics
	return diags
}

func resourceSettingsUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	var requirements api.PasswordRequirements

	for i := range d.Get("password_requirements").([]interface{}) {
		requirements.LengthMin = d.Get(fmt.Sprintf("password_requirements.%d.minimum_length", i)).(int)
		requirements.LowercaseMin = d.Get(fmt.Sprintf("password_requirements.%d.minimum_lowercase", i)).(int)
		requirements.UppercaseMin = d.Get(fmt.Sprintf("password_requirements.%d.minimum_uppercase", i)).(int)
		requirements.NumberMin = d.Get(fmt.Sprintf("password_requirements.%d.minimum_numbers", i)).(int)
		requirements.SymbolMin = d.Get(fmt.Sprintf("password_requirements.%d.minimum_symbols", i)).(int)
	}

	if _, err := client.UpdateSettings(ctx, &api.Settings{PasswordRequirements: requirements}); err != nil {
		return diag.FromErr(err)
	}

	return resourceSettingsRead(ctx, d, m)
}

func resourceSettingsDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	d.SetId("")
	var diags diag.Diagnostics
	return diags
}
