package provider

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
)

func dataSourceGroups() *schema.Resource {
	return &schema.Resource{
		Description: "Get a list of Infra groups.",

		ReadContext: dataSourceGroupsRead,

		Schema: map[string]*schema.Schema{
			"filter": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Description: "The name of the group.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"user_id": &schema.Schema{
							Description:      "The ID of the user who belongs to this group.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateStringIsID(),
							ConflictsWith: []string{
								"filter.0.user_name",
							},
						},
						"user_name": &schema.Schema{
							Description:      "The name of the user who belongs to this group.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateStringIsEmail(),
							ConflictsWith: []string{
								"filter.0.user_id",
							},
						},
					},
				},
			},
			"include_users": &schema.Schema{
				Description: "Include each group's members. Default is `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"groups": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Description: "The ID of the group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": &schema.Schema{
							Description: "The name of the group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"users": &schema.Schema{
							Description: "The members of the group.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func dataSourceGroupsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	request := api.ListGroupsRequest{
		PaginationRequest: api.PaginationRequest{
			Limit: 1000,
		},
	}

	for i := range d.Get("filter").([]interface{}) {
		request.Name = d.Get(fmt.Sprintf("filter.%d.name", i)).(string)

		if d.Get(fmt.Sprintf("filter.%d.user_id", i)).(string) != "" || d.Get(fmt.Sprintf("filter.%d.user_name", i)).(string) != "" {
			user, err := userFromIDOrEmail(ctx, client, d, fmt.Sprintf("filter.%d.user_id", i), fmt.Sprintf("filter.%d.user_name", i))
			if err != nil {
				return diag.FromErr(err)
			}

			request.UserID = user.ID
		}
	}

	response, err := client.ListGroups(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	sha1sum := sha1.New()

	groups := make([]map[string]interface{}, 0, response.Count)
	for _, item := range response.Items {
		group := make(map[string]interface{})
		group["id"] = item.ID.String()
		group["name"] = item.Name

		io.WriteString(sha1sum, item.ID.String())

		if d.Get("include_users").(bool) {
			request := api.ListUsersRequest{
				Group: item.ID,
				PaginationRequest: api.PaginationRequest{
					Limit: 1000,
				},
			}

			response, err := client.ListUsers(ctx, request)
			if err != nil {
				return diag.FromErr(err)
			}

			users := make([]string, 0, response.Count)
			for _, user := range response.Items {
				users = append(users, user.Name)
			}

			group["users"] = users
		}

		groups = append(groups, group)
	}

	if err := d.Set("groups", groups); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hex.EncodeToString(sha1sum.Sum(nil)))

	var diags diag.Diagnostics
	return diags
}
