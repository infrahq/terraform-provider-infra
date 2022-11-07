package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/infrahq/infra/api"
)

func resourceGroupRole() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an Infra group role. This resource can be used to assign an Infra role to groups.",

		CreateContext: resourceGroupRoleCreate,
		ReadContext:   resourceGroupRoleRead,
		DeleteContext: resourceGroupRoleDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The group role's unique identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"group_id": {
				Description:      "The ID of the group to assign this role.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsID(),
				ExactlyOneOf: []string{
					"group_id", "group_name",
				},
			},
			"group_name": {
				Description:      "The name of the group to assign this role.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsName(),
				ExactlyOneOf: []string{
					"group_id", "group_name",
				},
			},
			"role": {
				Description: "The name of the Infra role to assign to the group. Valid roles are `admin` or `view`",
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

func resourceGroupRoleCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	group, err := groupFromIDOrName(ctx, client, d, "group_id", "group_name")
	if err != nil {
		return diag.FromErr(err)
	}

	request := &api.GrantRequest{
		Group:     group.ID,
		Privilege: d.Get("role").(string),
		Resource:  "infra",
	}

	response, err := client.CreateGrant(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("group_id", group.ID.String()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.ID.String())
	return resourceGroupRoleRead(ctx, d, m)
}

func resourceGroupRoleRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	grant, err := client.GetGrant(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := client.GetGroup(ctx, grant.Group)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("group_name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceGroupRoleDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
