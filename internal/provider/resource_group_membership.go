package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
	"github.com/infrahq/infra/uid"
)

func resourceGroupMembership() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an Infra user grant. This resource can be used to assign groups to users.",

		CreateContext: resourceGroupMembershipCreate,
		ReadContext:   resourceGroupMembershipRead,
		DeleteContext: resourceGroupMembershipDelete,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Description:      "The ID of the user to assign to the group.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsID(),
				ExactlyOneOf: []string{
					"user_id", "user_name",
				},
			},
			"user_name": {
				Description:      "The email of the user to assign to the group.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsEmail(),
				ExactlyOneOf: []string{
					"user_id", "user_name",
				},
			},
			"group_id": {
				Description:      "The ID of the group to assign to the user.",
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
				Description: "The name of the group to assign to the user.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				ExactlyOneOf: []string{
					"group_id", "group_name",
				},
			},
		},
	}
}

func resourceGroupMembershipCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	var diags diag.Diagnostics

	user, err := userFromIDOrEmail(ctx, client, d, "user_id", "user_name")
	if err != nil {
		diags = append(diags, diag.Diagnostic{Summary: err.Error()})
	}

	group, err := groupFromIDOrName(ctx, client, d, "group_id", "group_name")
	if err != nil {
		diags = append(diags, diag.Diagnostic{Summary: err.Error()})
	}

	if diags.HasError() {
		return diags
	}

	request := &api.UpdateUsersInGroupRequest{
		GroupID:      group.ID,
		UserIDsToAdd: []uid.ID{user.ID},
	}

	if err := client.UpdateUsersInGroup(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("user_id", user.ID.String()); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("group_id", group.ID.String()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", user.Name, group.Name))
	return resourceGroupMembershipRead(ctx, d, m)
}

func resourceGroupMembershipRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	userID, err := ParseID(d, "user_id")
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := client.GetUser(ctx, userID)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := ParseID(d, "group_id")
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := client.GetGroup(ctx, groupID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("user_name", user.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("group_name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceGroupMembershipDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	userID, err := ParseID(d, "user_id")
	if err != nil {
		return diag.FromErr(err)
	}

	groupID, err := ParseID(d, "group_id")
	if err != nil {
		return diag.FromErr(err)
	}

	request := &api.UpdateUsersInGroupRequest{
		GroupID:         groupID,
		UserIDsToRemove: []uid.ID{userID},
	}

	if err := client.UpdateUsersInGroup(ctx, request); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	var diags diag.Diagnostics
	return diags
}
