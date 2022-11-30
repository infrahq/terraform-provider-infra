package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/infrahq/infra/api"
)

func resourceIdentityProvider() *schema.Resource {
	return &schema.Resource{
		Description: "The infra identity provider resource can be used to create and manage identity providers.",

		CreateContext: resourceIdentityProviderCreate,
		ReadContext:   resourceIdentityProviderRead,
		UpdateContext: resourceIdentityProviderUpdate,
		DeleteContext: resourceIdentityProviderDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The identity provider's unique identifier.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description:      "The identity provider's name. If omitted, a name will be automatically generated. Identity provider names may include letters (uppercase and lowercase), numbers, underscores `_`, hyphens `-`, and periods `.`.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validateStringIsName(),
			},
			"issuer": {
				Description:      "The identity provider's full authorization server URL. Must start with `https://`.",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
				ExactlyOneOf: []string{
					"issuer", "google", "azure", "okta",
				},
			},
			"client_id": {
				Description: "The identity provider's OIDC client ID.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_secret": {
				Description: "The identity provider's OIDC client secret.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"azure": {
				Description: "Azure AD identity provider configurations.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ExactlyOneOf: []string{
					"issuer", "google", "azure", "okta",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Description: "The Azure AD tenant ID.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"google": {
				Description: "Google identity provider configurations.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ExactlyOneOf: []string{
					"issuer", "google", "azure", "okta",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_email": {
							Description:      "A Google workspace admin user email. Infra will impersonate this user when making API calls to retrieve Google groups. If set, `service_account_key` must also be set.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: validateStringIsEmail(),
							RequiredWith: []string{
								"google.0.admin_email", "google.0.service_account_key",
							},
						},
						"service_account_key": {
							Description:      "A Google service account private key file. Must be a JSON-formatted string. If set, `admin_email` must also be set.",
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
							RequiredWith: []string{
								"google.0.admin_email", "google.0.service_account_key",
							},
						},
					},
				},
			},
			"okta": {
				Description: "Okta identity provider configurations.",
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				ExactlyOneOf: []string{
					"issuer", "google", "azure", "okta",
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer": {
							Description:      "The full Okta authorization server URL. Must start with `https://`.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
						},
					},
				},
			},
		},
	}
}

func resourceIdentityProviderCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	request := &api.CreateProviderRequest{
		Name:         d.Get("name").(string),
		URL:          d.Get("issuer").(string),
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
	}

	for i := range d.Get("azure").([]interface{}) {
		request.Kind = "azure"
		request.URL = fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", d.Get(fmt.Sprintf("azure.%d.tenant_id", i)).(string))
		break
	}

	for i := range d.Get("google").([]interface{}) {
		request.API = &api.ProviderAPICredentials{
			DomainAdminEmail: d.Get(fmt.Sprintf("google.%d.admin_email", i)).(string),
		}

		if serviceAccountKey := d.Get(fmt.Sprintf("google.%d.service_account_key", i)).(string); serviceAccountKey != "" {
			var serviceAccountKeyFile struct {
				PrivateKey  string `json:"service_account_key"`
				ClientEmail string `json:"client_email"`
			}

			if err := json.Unmarshal([]byte(serviceAccountKey), &serviceAccountKeyFile); err != nil {
				return diag.FromErr(err)
			}

			request.API.PrivateKey = api.PEM(serviceAccountKeyFile.PrivateKey)
			request.API.ClientEmail = serviceAccountKeyFile.ClientEmail
		}

		request.Kind = "google"
		request.URL = "https://accounts.google.com"
		break
	}

	for i := range d.Get("okta").([]interface{}) {
		request.Kind = "okta"
		request.URL = d.Get(fmt.Sprintf("okta.%d.issuer", i)).(string)
	}

	provider, err := client.CreateProvider(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(provider.ID.String())
	return resourceIdentityProviderRead(ctx, d, m)
}

func resourceIdentityProviderRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	provider, err := client.GetProvider(ctx, id)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", provider.Name); err != nil {
		return diag.FromErr(err)
	}

	providerURL := provider.URL
	if !strings.HasPrefix(providerURL, "https://") {
		providerURL = fmt.Sprintf("https://%s", providerURL)
	}

	if err := d.Set("issuer", providerURL); err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return diags
}

func resourceIdentityProviderUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	request := api.UpdateProviderRequest{
		ID:           id,
		Name:         d.Get("name").(string),
		URL:          d.Get("issuer").(string),
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
	}

	for i := range d.Get("azure").([]interface{}) {
		request.Kind = "azure"
		request.URL = fmt.Sprintf("https://login.microsoftonline.com/%s/v2.0", d.Get(fmt.Sprintf("azure.%d.tenant_id", i)).(string))
		break
	}

	for i := range d.Get("google").([]interface{}) {
		request.API = &api.ProviderAPICredentials{
			DomainAdminEmail: d.Get(fmt.Sprintf("google.%d.admin_email", i)).(string),
		}

		if serviceAccountKey := d.Get(fmt.Sprintf("google.%d.service_account_key", i)).(string); serviceAccountKey != "" {
			var serviceAccountKeyFile struct {
				PrivateKey  string `json:"service_account_key"`
				ClientEmail string `json:"client_email"`
			}

			if err := json.Unmarshal([]byte(serviceAccountKey), &serviceAccountKeyFile); err != nil {
				return diag.FromErr(err)
			}

			request.API.PrivateKey = api.PEM(serviceAccountKeyFile.PrivateKey)
			request.API.ClientEmail = serviceAccountKeyFile.ClientEmail
		}

		request.Kind = "google"
		request.URL = "https://accounts.google.com"
		break
	}

	for i := range d.Get("okta").([]interface{}) {
		request.Kind = "okta"
		request.URL = d.Get(fmt.Sprintf("okta.%d.issuer", i)).(string)
	}

	_, err = client.UpdateProvider(ctx, request)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIdentityProviderRead(ctx, d, m)
}

func resourceIdentityProviderDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*api.Client)

	id, err := ParseID(d, "id")
	if err != nil {
		return diag.FromErr(err)
	}

	if err := client.DeleteProvider(ctx, id); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	var diags diag.Diagnostics
	return diags
}
