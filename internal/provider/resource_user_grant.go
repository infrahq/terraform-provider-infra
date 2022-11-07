package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
)

func resourceUserGrant() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an Infra user grant. This resource can be used to assign grants to users.",

		CreateContext: resourceUserGrantCreate,
		ReadContext:   resourceUserGrantRead,
		DeleteContext: resourceUserGrantDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The grant's unique identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_id": {
				Description:      "The ID of the user to assign this grant.",
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
				Description:      "The email of the user to assign this grant.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsEmail(),
				ExactlyOneOf: []string{
					"user_id", "user_email",
				},
			},
			"kubernetes": {
				Description: "Kubernetes user grant configurations.",
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				ExactlyOneOf: []string{
					"kubernetes",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Description: "The name of the Kubernetes ClusterRole to assign to the user.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
						"cluster": {
							Description: "The name of the Kubernetes cluster to assign to the user.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
						},
						"namespace": {
							Description: "The namespace of the Kubernetes cluster to assign to the name.",
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},
		},
	}
}

func resourceUserGrantCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	user, err := userFromIDOrEmail(ctx, client, d, "user_id", "user_email")
	if err != nil {
		return diag.FromErr(err)
	}

	request := &api.GrantRequest{
		User: user.ID,
	}

	for i := range d.Get("kubernetes").([]interface{}) {
		resource := d.Get(fmt.Sprintf("kubernetes.%d.cluster", i)).(string)
		if namespace := d.Get(fmt.Sprintf("kubernetes.%d.namespace", i)).(string); namespace != "" {
			resource = fmt.Sprintf("%s.%s", resource, namespace)
		}

		request.Resource = resource
		request.Privilege = d.Get(fmt.Sprintf("kubernetes.%d.role", i)).(string)
		break
	}

	response, err := client.CreateGrant(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("user_id", user.ID.String()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.ID.String())
	return resourceUserGrantRead(ctx, d, m)
}

func resourceUserGrantRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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

func resourceUserGrantDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
