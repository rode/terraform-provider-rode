package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/common"
	"google.golang.org/grpc"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Description: "",
				Type:        schema.TypeString,
				Required:    true,
			},
			"disable_transport_security": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			// TODO: separate oidc/basic objects instead?
			"oidc_client_id": {
				Description:   "",
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"oidc_client_secret"},
				ConflictsWith: []string{"username", "password"},
			},
			"oidc_client_secret": {
				Description:   "",
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"oidc_client_id"},
				ConflictsWith: []string{"username", "password"},
				Sensitive:     true,
			},
			"oidc_token_url": {
				Description:  "",
				Optional:     true,
				Type:         schema.TypeString,
				RequiredWith: []string{"oidc_client_id"},
			},
			"oidc_scopes": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
			},
			"oidc_tls_insecure_skip_verify": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"username": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"password"},
				ConflictsWith: []string{"oidc_client_id", "oidc_client_secret"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				RequiredWith:  []string{"username"},
				ConflictsWith: []string{"oidc_client_id", "oidc_client_secret"},
				Sensitive:     true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rode_policy_group": resourcePolicyGroup(),
			"rode_policy": resourcePolicy(),
			"rode_assignment": resourcePolicyAssignment(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
		// TODO: configurable lazy init of client

		config := &common.ClientConfig{
			Rode: &common.RodeClientConfig{
				Host: data.Get("host").(string),
				DisableTransportSecurity: data.Get("disable_transport_security").(bool),
			},
			OIDCAuth: &common.OIDCAuthConfig{
				ClientID:     data.Get("oidc_client_id").(string),
				ClientSecret: data.Get("oidc_client_secret").(string),
				TokenURL:     data.Get("oidc_token_url").(string),
				TlsInsecureSkipVerify: data.Get("oidc_tls_insecure_skip_verify").(bool),
				Scopes:       data.Get("oidc_scopes").(string),
			},
			BasicAuth: &common.BasicAuthConfig{
				Username: data.Get("username").(string),
				Password: data.Get("password").(string),
			},
		}

		rode, err := common.NewRodeClient(
			config,
			// TODO: source version, either by embedding file or with goreleaser
			grpc.WithUserAgent(provider.UserAgent("terraform-provider-rode", "0.0.1")),
		)

		return rode, diag.FromErr(err)
	}

	return provider
}
