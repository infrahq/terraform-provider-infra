package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/infrahq/infra/api"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		var sb strings.Builder

		sb.WriteString(s.Description)

		if s.ConflictsWith != nil && len(s.ConflictsWith) > 0 {
			fmt.Fprintf(&sb, " Cannot be used with `%s`.", strings.Join(s.ConflictsWith, "`, `"))
		}

		if s.ExactlyOneOf != nil && len(s.ExactlyOneOf) > 0 {
			fmt.Fprintf(&sb, " One of `%s` must be set.", strings.Join(s.ExactlyOneOf, "`, `"))
		}

		if s.Default != nil {
			fmt.Fprintf(&sb, " Default is %q.", s.Default)
		}

		return strings.TrimSpace(sb.String())
	}
}

func New() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: configure(),

		Schema: map[string]*schema.Schema{
			"host": &schema.Schema{
				Description:      "The Infra server instance Terraform will communicate with. Can also be sourced from `INFRA_HOST`. Default is `https://api.infrahq.com`.",
				Type:             schema.TypeString,
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc("INFRA_HOST", "https://api.infrahq.com"),
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
			},
			"access_key": &schema.Schema{
				Description: "The access key used to authenticate with the Infra server. Can also be sourced from `INFRA_ACCESS_KEY`.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("INFRA_ACCESS_KEY", nil),
			},
			"skip_tls_verify": &schema.Schema{
				Description: "Controls client verification of the server certificate. This should only be `true` for testing or development. Can also be sourced from`INFRA_SKIP_TLS_VERIFY`.",
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INFRA_SKIP_TLS_VERIFY", nil),
				ConflictsWith: []string{
					"server_certificate",
					"server_certificate_file",
				},
			},
			"server_certificate": &schema.Schema{
				Description:      "The server's PEM-encoded public certificate for client verification. Can also be sourced from `INFRA_SERVER_CERTIFICATE`.",
				Type:             schema.TypeString,
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc("INFRA_SERVER_CERTIFICATE", nil),
				ValidateDiagFunc: validateStringIsPEMEncoded(),
				ConflictsWith: []string{
					"skip_tls_verify",
					"server_certificate_file",
				},
			},
			"server_certificate_file": &schema.Schema{
				Description:      "The server's PEM-encoded public certificate file for client verification. Can also be sourced from `INFRA_SERVER_CERTIFICATE_FILE`.",
				Type:             schema.TypeString,
				Optional:         true,
				DefaultFunc:      schema.EnvDefaultFunc("INFRA_SERVER_CERTIFICATE_FILE", nil),
				ValidateDiagFunc: validateStringIsPEMEncodedFile(),
				ConflictsWith: []string{
					"skip_tls_verify",
					"server_certificate",
				},
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"infra_destinations": dataSourceDestinations(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"infra_user":              resourceUser(),
			"infra_group":             resourceGroup(),
			"infra_group_membership":  resourceGroupMembership(),
			"infra_grant":             resourceGrant(),
			"infra_identity_provider": resourceIdentityProvider(),
		},
	}
}

func configure() func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		pool, err := x509.SystemCertPool()
		if err != nil {
			return nil, diag.FromErr(err)
		}

		if cacert := d.Get("server_certificate").(string); cacert != "" {
			b, err := DecodePEM([]byte(cacert), "CERTIFICATE")
			if err != nil {
				return nil, diag.FromErr(err)
			}

			if ok := pool.AppendCertsFromPEM(b); !ok {
				return nil, diag.Errorf("not ok %#v", string(b))
			}
		}

		if cacertfile := d.Get("server_certificate_file").(string); cacertfile != "" {
			b, err := DecodePEMFile(cacertfile, "CERTIFICATE")
			if err != nil {
				return nil, diag.FromErr(err)
			}

			if ok := pool.AppendCertsFromPEM(b); !ok {
				return nil, diag.Errorf("not ok %#v", string(b))
			}
		}

		return &api.Client{
			Name:      "terraform",
			Version:   "0.17.1",
			URL:       d.Get("host").(string),
			AccessKey: d.Get("access_key").(string),
			HTTP: http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: d.Get("skip_tls_verify").(bool),
						RootCAs:            pool,
					},
				},
			},
		}, nil
	}
}
