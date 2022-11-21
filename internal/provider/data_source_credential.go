package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
)

func dataSourceCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Get an authentication token to communicate with a registered destination.",

		ReadContext: dataSourceCredentialRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"token": &schema.Schema{
				Description: "A token that can be used to authenticate the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceCredentialRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	credential, err := client.CreateToken(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("token", credential.Token); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	var diags diag.Diagnostics
	return diags
}
