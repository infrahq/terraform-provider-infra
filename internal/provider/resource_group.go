package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
	"github.com/infrahq/infra/uid"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an Infra group. This resource can be used to create and manage groups.",

		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		DeleteContext: resourceGroupDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The group's unique identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:      "The group's name. Group names may include letters (uppercase and lowercase), numbers, underscores `_`, hyphens `-`, and periods `.`.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsName(),
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	name := strings.TrimSpace(d.Get("name").(string))
	group, err := client.CreateGroup(ctx, &api.CreateGroupRequest{Name: name})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.ID.String())
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	group, err := client.GetGroup(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", group.Name); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.DeleteGroup(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	var diags diag.Diagnostics
	return diags
}

func groupFromIDOrName(ctx context.Context, client *api.Client, d *schema.ResourceData, id, name string) (*api.Group, error) {
	if s := d.Get(id).(string); s != "" {
		groupID, err := uid.Parse([]byte(s))
		if err != nil {
			return nil, err
		}

		return client.GetGroup(ctx, groupID)
	}

	if s := d.Get(name).(string); s != "" {
		request := api.ListGroupsRequest{
			Name: d.Get(name).(string),
			PaginationRequest: api.PaginationRequest{
				Limit: 1,
			},
		}

		response, err := client.ListGroups(ctx, request)
		if err != nil {
			return nil, err
		}

		if response.Count < 1 {
			return nil, fmt.Errorf("group not found")
		}

		return &response.Items[0], nil
	}

	return nil, fmt.Errorf("one of `%s,%s` must be specified", id, name)
}
