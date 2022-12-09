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

func dataSourceDestinations() *schema.Resource {
	return &schema.Resource{
		Description: "Get a list of Infra destinations.",

		ReadContext: dataSourceDestinationsRead,

		Schema: map[string]*schema.Schema{
			"filter": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Description: "The name of the destination.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"kind": &schema.Schema{
							Description: "The kind of the destination.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"destinations": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Description: "The ID of the destination.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": &schema.Schema{
							Description: "The name of the destination.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"kind": &schema.Schema{
							Description: "The kind of the destination.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDestinationsRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	request := api.ListDestinationsRequest{
		PaginationRequest: api.PaginationRequest{
			Limit: 1000,
		},
	}

	for i := range d.Get("filter").([]interface{}) {
		request.Name = d.Get(fmt.Sprintf("filter.%d.name", i)).(string)
		request.Kind = d.Get(fmt.Sprintf("filter.%d.kind", i)).(string)
	}

	response, err := client.ListDestinations(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	sha1sum := sha1.New()

	destinations := make([]map[string]interface{}, 0, response.Count)
	for _, item := range response.Items {
		destination := make(map[string]interface{})
		destination["id"] = item.ID.String()
		destination["name"] = item.Name
		destination["kind"] = item.Kind

		io.WriteString(sha1sum, item.ID.String())

		destinations = append(destinations, destination)
	}

	if err := d.Set("destinations", destinations); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hex.EncodeToString(sha1sum.Sum(nil)))

	var diags diag.Diagnostics
	return diags
}
