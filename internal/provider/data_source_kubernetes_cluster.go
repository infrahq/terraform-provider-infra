package provider

import (
	"context"
	"encoding/base64"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
)

func dataSourceKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Get information about a registered Kubernetes cluster.",

		ReadContext: dataSourceKubernetesClusterRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Description: "The Kubernetes cluster's unique identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": &schema.Schema{
				Description: "The Kubernetes cluster's name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"endpoint": &schema.Schema{
				Description: "The Kubernetes cluster's API server connection endpoint.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"certificate_authority": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data": &schema.Schema{
							Description: "The Kubernetes cluster's base64-encoded certificate authority data. Set this as `certificate-authority-data` in `kubeconfig` for this cluster.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"namespaces": &schema.Schema{
				Description: "A list of known namespaces for this Kubernetes cluster.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"roles": &schema.Schema{
				Description: "A list of known ClusterRoles for this Kubernetes cluster. To make a ClusterRole known to Infra, add `app.infrahq.com/include-role='true'` as a label value to the ClusterRole.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	name := d.Get("name").(string)

	request := api.ListDestinationsRequest{
		Name: name,
		PaginationRequest: api.PaginationRequest{
			Limit: 1,
		},
	}

	response, err := client.ListDestinations(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	if response.Count < 1 {
		return diag.Errorf("%s not found", name)
	}

	cluster := response.Items[0]

	if err := d.Set("name", cluster.Name); err != nil {
		return diag.FromErr(err)
	}

	endpoint := url.URL{
		Scheme: "https",
		Host:   cluster.Connection.URL,
	}

	if err := d.Set("endpoint", endpoint.String()); err != nil {
		return diag.FromErr(err)
	}

	caData := []map[string]string{
		map[string]string{
			"data": base64.StdEncoding.EncodeToString([]byte(cluster.Connection.CA)),
		},
	}

	if err := d.Set("certificate_authority", caData); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("namespaces", cluster.Resources); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("roles", cluster.Roles); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster.ID.String())

	var diags diag.Diagnostics
	return diags
}
