package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/infrahq/infra/api"
)

func resourceGrant() *schema.Resource {
	return &schema.Resource{
		Description: "Provides an Infra grant. This resource can be used to assign grants to users or groups.",

		CreateContext: resourceGrantCreate,
		ReadContext:   resourceGrantRead,
		DeleteContext: resourceGrantDelete,

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
					"user_id", "user_name", "group_id", "group_name",
				},
			},
			"user_name": {
				Description:      "The email of the user to assign this grant.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsEmail(),
				ExactlyOneOf: []string{
					"user_id", "user_name", "group_id", "group_name",
				},
			},
			"group_id": {
				Description:      "The ID of the group to assign this grant.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
				ValidateDiagFunc: validateStringIsID(),
				ExactlyOneOf: []string{
					"user_id", "user_name", "group_id", "group_name",
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
					"user_id", "user_name", "group_id", "group_name",
				},
			},
			"infra": {
				Description: "Infra grant configurations.",
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				ExactlyOneOf: []string{
					"infra", "kubernetes",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role": {
							Description: "The name of the Infra role to assign to the user. Valid roles are `admin` or `view`.",
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							ValidateDiagFunc: validation.ToDiagFunc(
								validation.StringInSlice([]string{"admin", "view"}, true),
							),
						},
					},
				},
			},
			"kubernetes": {
				Description: "Kubernetes grant configurations.",
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				ExactlyOneOf: []string{
					"infra", "kubernetes",
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

func resourceGrantCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	request := &api.GrantRequest{}

	if d.Get("user_id").(string) != "" || d.Get("user_name").(string) != "" {
		user, err := userFromIDOrEmail(ctx, client, d, "user_id", "user_name")
		if err != nil {
			return diag.FromErr(err)
		}

		request.User = user.ID
	}

	if d.Get("group_id").(string) != "" || d.Get("group_name").(string) != "" {
		group, err := groupFromIDOrName(ctx, client, d, "group_id", "group_name")
		if err != nil {
			return diag.FromErr(err)
		}

		request.Group = group.ID
	}

	for i := range d.Get("infra").([]interface{}) {
		request.Resource = "infra"
		request.Privilege = d.Get(fmt.Sprintf("infra.%d.role", i)).(string)
		break
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

	d.SetId(response.ID.String())
	return resourceGrantRead(ctx, d, m)
}

func resourceGrantRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	grant, err := client.GetGrant(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if grant.User != 0 {
		user, err := client.GetUser(ctx, grant.User)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("user_id", user.ID.String()); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("user_name", user.Name); err != nil {
			return diag.FromErr(err)
		}
	}

	if grant.Group != 0 {
		group, err := client.GetGroup(ctx, grant.Group)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("group_id", group.ID.String()); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("group_name", group.Name); err != nil {
			return diag.FromErr(err)
		}
	}

	var diags diag.Diagnostics
	return diags
}

func resourceGrantDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
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
