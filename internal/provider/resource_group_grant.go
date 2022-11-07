package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/infrahq/infra/api"
)

func resourceGroupGrant() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an Infra group grant. This resource can be used to assign grants to groups.",

		CreateContext: resourceGroupGrantCreate,
		ReadContext:   resourceGroupGrantRead,
		DeleteContext: resourceGroupGrantDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The grant's unique identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"group_id": {
				Description:      "The ID of the group to assign this grant.",
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
				Description:      "The name of the group to assign this grant.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsName(),
				ExactlyOneOf: []string{
					"group_id", "group_name",
				},
			},
			"kubernetes": {
				Description: "Kubernetes group grant configurations.",
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
							Description: "The name of the Kubernetes ClusterRole to assign to the group.",
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

func resourceGroupGrantCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	group, err := groupFromIDOrName(ctx, client, d, "group_id", "group_name")
	if err != nil {
		return diag.FromErr(err)
	}

	request := &api.GrantRequest{
		Group: group.ID,
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

	if err := d.Set("group_id", group.ID.String()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.ID.String())
	return resourceGroupGrantRead(ctx, d, m)
}

func resourceGroupGrantRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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

func resourceGroupGrantDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
