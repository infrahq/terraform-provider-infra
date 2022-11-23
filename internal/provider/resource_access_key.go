package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/infrahq/infra/api"
)

func resourceAccessKey() *schema.Resource {
	return &schema.Resource{
		Description: `Provides an Infra access key. This resource can be used to create and manage connector access keys.

-> This resource requires Infra server version 0.20.0 or higher.`,

		CreateContext: resourceAccessKeyCreate,
		ReadContext:   resourceAccessKeyRead,
		DeleteContext: resourceAccessKeyDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The access key's unique identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:      "The access key's name. If omitted, a name will be automatically generated. Identity provider names may include letters (uppercase and lowercase), numbers, underscores `_`, hyphens `-`, and periods `.`.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsName(),
			},
			"secret": {
				Description: "The access key secret.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"expires_in": {
				Description:      `The total amount of time before the access key expires. Format is a duration string, a sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300s" or "2h45m". Valid time units are "s", "m", "h". Default is 8766h0m0s.`,
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsDuration(),
				DiffSuppressFunc: DurationDiffSuppressFunc(),
				ConflictsWith: []string{
					"expires_at",
				},
			},
			"expires_at": {
				Description:      `The date-time when the access key will expire. Format is a RFC3339 timestamp, e.g. "2006-01-02T15:04:05Z07:00."`,
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
				ConflictsWith: []string{
					"expires_in",
				},
			},
			"inactivity_timeout": {
				Description:      `The amount of time before the access key expires if left unused. If the access key is used before it expires, it will be renewed for the same duration. Format is a duration string, a sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300s" or "2h45m". Valid time units are "s", "m", "h". If value is greater than or equal to the remaining lifetime of the access key, the access key will not be renewed. Default is 72h0m0s.`,
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsDuration(),
				DiffSuppressFunc: DurationDiffSuppressFunc(),
			},
		},
	}
}

func resourceAccessKeyCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	if err := requireMinimumServerVersion(ctx, client, "0.20.0"); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	var user *api.User
	var err error

	user, err = userFromEmail(ctx, client, "connector")
	if err != nil {
		diags = append(diags, diag.Diagnostic{Summary: err.Error()})
	}

	expires, err := ParseDuration(d, "expires_in", "expires_at", "8766h0m0s")
	if err != nil {
		diags = append(diags, diag.Diagnostic{Summary: err.Error()})
	}

	inactivity, err := ParseDuration(d, "inactivity_timeout", "", expires.String())
	if err != nil {
		diags = append(diags, diag.Diagnostic{Summary: err.Error()})
	}

	if inactivity > expires {
		inactivity = expires
	}

	if diags.HasError() {
		return diags
	}

	request := &api.CreateAccessKeyRequest{
		UserID:            user.ID,
		Name:              d.Get("name").(string),
		InactivityTimeout: api.Duration(inactivity),
		Expiry:            api.Duration(expires),
	}

	response, err := client.CreateAccessKey(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", response.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("secret", response.AccessKey); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("expires_in", expires.Truncate(time.Second).String()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.ID.String())
	return resourceAccessKeyRead(ctx, d, m)
}

func resourceAccessKeyRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	request := api.ListAccessKeysRequest{
		Name:        d.Get("name").(string),
		ShowExpired: true,
		PaginationRequest: api.PaginationRequest{
			Limit: 1,
		},
	}

	response, err := client.ListAccessKeys(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.Count < 1 {
		return diag.Errorf("access key not found")
	}

	accessKey := response.Items[0]

	if err := d.Set("expires_at", accessKey.Expires.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceAccessKeyDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.DeleteAccessKey(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	var diags diag.Diagnostics
	return diags
}

func ParseDuration(d *schema.ResourceData, duration, timestamp, defaultValue string) (time.Duration, error) {
	if duration != "" {
		if s := d.Get(duration).(string); s != "" {
			dur, err := time.ParseDuration(s)
			if err != nil {
				return time.Duration(0), err
			}

			return dur, nil
		}
	}

	if timestamp != "" {
		if s := d.Get(timestamp).(string); s != "" {
			ts, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return time.Duration(0), err
			}

			return time.Until(ts), nil
		}
	}

	return time.ParseDuration(defaultValue)
}
