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

func dataSourceUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Get a list of Infra users.",

		ReadContext: dataSourceUsersRead,

		Schema: map[string]*schema.Schema{
			"filter": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Description:      "The name of the user.",
							Type:             schema.TypeString,
							ValidateDiagFunc: validateStringIsEmail(),
							Optional:         true,
						},
						"group_id": &schema.Schema{
							Description:      "The ID of the group where user is a member.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateStringIsID(),
							ConflictsWith: []string{
								"filter.0.group_name",
							},
						},
						"group_name": &schema.Schema{
							Description: "The name of the group where user is a member.",
							Type:        schema.TypeString,
							Optional:    true,
							ConflictsWith: []string{
								"filter.0.group_id",
							},
						},
					},
				},
			},
			"include_groups": &schema.Schema{
				Description: "Include each user's group membership. Default is `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"users": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Description: "The ID of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": &schema.Schema{
							Description: "The name of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"groups": &schema.Schema{
							Description: "The groups the user is a member of.",
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

func dataSourceUsersRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	request := api.ListUsersRequest{
		PaginationRequest: api.PaginationRequest{
			Limit: 1000,
		},
	}

	for i := range d.Get("filter").([]interface{}) {
		request.Name = d.Get(fmt.Sprintf("filter.%d.name", i)).(string)

		if d.Get(fmt.Sprintf("filter.%d.group_id", i)).(string) != "" || d.Get(fmt.Sprintf("filter.%d.group_name", i)).(string) != "" {
			group, err := groupFromIDOrName(ctx, client, d, fmt.Sprintf("filter.%d.group_id", i), fmt.Sprintf("filter.%d.group_name", i))
			if err != nil {
				return diag.FromErr(err)
			}

			request.Group = group.ID
		}
	}

	response, err := client.ListUsers(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	sha1sum := sha1.New()

	users := make([]map[string]interface{}, 0, response.Count)
	for _, item := range response.Items {
		user := make(map[string]interface{})
		user["id"] = item.ID.String()
		user["name"] = item.Name

		io.WriteString(sha1sum, item.ID.String())

		if d.Get("include_groups").(bool) {
			request := api.ListGroupsRequest{
				UserID: item.ID,
				PaginationRequest: api.PaginationRequest{
					Limit: 1000,
				},
			}

			response, err := client.ListGroups(ctx, request)
			if err != nil {
				return diag.FromErr(err)
			}

			groups := make([]string, 0, response.Count)
			for _, group := range response.Items {
				groups = append(groups, group.Name)
			}

			user["groups"] = groups
		}

		users = append(users, user)
	}

	if err := d.Set("users", users); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hex.EncodeToString(sha1sum.Sum(nil)))

	var diags diag.Diagnostics
	return diags
}
