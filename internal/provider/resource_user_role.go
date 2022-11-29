package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/infrahq/infra/api"
)

func resourceUserRole() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an Infra user role. This resource can be used to assign an Infra role to users.",

		CreateContext: resourceUserRoleCreate,
		ReadContext:   resourceUserRoleRead,
		DeleteContext: resourceUserRoleDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The user role's unique identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_id": {
				Description:      "The ID of the user to assign this role.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsID(),
				ExactlyOneOf: []string{
					"user_id", "user_email",
				},
			},
			"user_email": {
				Description:      "The email of the user to assign this role.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsEmail(),
				ExactlyOneOf: []string{
					"user_id", "user_email",
				},
			},
			"role": {
				Description: "The name of the Infra role to assign to the user. Valid roles are `admin` or `view`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice([]string{"admin", "view"}, true),
				),
			},
		},
	}
}

func resourceUserRoleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	user, err := userFromIDOrEmail(ctx, client, d, "user_id", "user_email")
	if err != nil {
		return diag.FromErr(err)
	}

	request := &api.GrantRequest{
		User:      user.ID,
		Privilege: d.Get("role").(string),
		Resource:  "infra",
	}

	response, err := client.CreateGrant(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("user_id", user.ID.String()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.ID.String())
	return resourceUserRoleRead(ctx, d, m)
}

func resourceUserRoleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	grant, err := client.GetGrant(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := client.GetUser(ctx, grant.User)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("user_email", user.Name); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceUserRoleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.DeleteGrant(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	var diags diag.Diagnostics
	return diags
}
