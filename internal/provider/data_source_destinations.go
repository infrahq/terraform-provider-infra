package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
)

func dataSourceDestinations() *schema.Resource {
	return &schema.Resource{
		Description: "Get a list of registered Infra destinations.",

		ReadContext: dataSourceDestinationsRead,

		Schema: map[string]*schema.Schema{
			"names": &schema.Schema{
				Description: "A list of registered destinations. Use `data.infra_destination` to retrieve information about individual destinations.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

	response, err := client.ListDestinations(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	names := make([]string, 0, response.Count)
	for _, cluster := range response.Items {
		names = append(names, cluster.Name)
	}

	if err := d.Set("names", names); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	var diags diag.Diagnostics
	return diags
}
