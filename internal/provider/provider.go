package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/common"
	"google.golang.org/grpc"
	"os"
	"strconv"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		provider := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"host": {
					Description: "Host and port of the Rode instance. Can also be specified by setting the `RODE_HOST` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"disable_transport_security": {
					Description: "Disables transport security for the gRPC connection to Rode. Can also be set with the `RODE_DISABLE_TRANSPORT_SECURITY` environment variable.",
					Type:        schema.TypeBool,
					Default:     false,
					Optional:    true,
				},
				"oidc_client_id": {
					Description:   "OIDC/OAuth2 client id that is permitted the client credentials grant. Can be set with the `RODE_OIDC_CLIENT_ID` environment variable.",
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{"basic_username", "basic_password"},
				},
				"oidc_client_secret": {
					Description:   "Corresponding client secret for oidc_client_id. Can be set with the `RODE_OIDC_CLIENT_SECRET` environment variable.",
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{"basic_username", "basic_password"},
					Sensitive:     true,
				},
				"oidc_token_url": {
					Description: "OAuth2 token url. Can be set with the OIDC_TOKEN_URL environment variable",
					Optional:    true,
					Type:        schema.TypeString,
				},
				"oidc_scopes": {
					Description: "A space-delimited list of scopes to request in the client credentials grant. Can also be set with the `RODE_OIDC_SCOPES` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
				},
				"oidc_tls_insecure_skip_verify": {
					Description: "Disable transport security when communicating with the OAuth2 server. Only recommended for local development. Set with the `RODE_OIDC_TLS_INSECURE_SKIP_VERIFY` environment variable.",
					Type:        schema.TypeBool,
					Optional:    true,
				},
				"basic_username": {
					Description:   "The username configured in the Rode instance for basic auth. Cannot be configured alongside any of the OIDC options. Can be set with the `RODE_BASIC_USERNAME` environment variable.",
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{"oidc_client_id", "oidc_client_secret"},
				},
				"basic_password": {
					Description:   "Corresponding password for basic_username. Can be set with the `RODE_BASIC_PASSWORD` environment variable.",
					Type:          schema.TypeString,
					Optional:      true,
					ConflictsWith: []string{"oidc_client_id", "oidc_client_secret"},
					Sensitive:     true,
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"rode_policy_group":      resourcePolicyGroup(),
				"rode_policy":            resourcePolicy(),
				"rode_policy_assignment": resourcePolicyAssignment(),
			},
		}

		provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			// TODO: configurable lazy init of client

			rodeDisableTransportSecurity, err := getProviderOptionBool(d, "disable_transport_security", "RODE_DISABLE_TRANSPORT_SECURITY")
			if err != nil {
				return nil, diag.FromErr(err)
			}
			oidcInsecureSkipVerify, err := getProviderOptionBool(d, "oidc_tls_insecure_skip_verify", "RODE_OIDC_TLS_INSECURE_SKIP_VERIFY")
			if err != nil {
				return nil, diag.FromErr(err)
			}

			config := &common.ClientConfig{
				Rode: &common.RodeClientConfig{
					Host:                     getProviderOption(d, "host", "RODE_HOST"),
					DisableTransportSecurity: rodeDisableTransportSecurity,
				},
				OIDCAuth: &common.OIDCAuthConfig{
					ClientID:              getProviderOption(d, "oidc_client_id", "RODE_OIDC_CLIENT_ID"),
					ClientSecret:          getProviderOption(d, "oidc_client_secret", "RODE_OIDC_CLIENT_SECRET"),
					TokenURL:              getProviderOption(d, "oidc_token_url", "RODE_OIDC_TOKEN_URL"),
					TlsInsecureSkipVerify: oidcInsecureSkipVerify,
					Scopes:                getProviderOption(d, "oidc_scopes", "RODE_OIDC_SCOPES"),
				},
				BasicAuth: &common.BasicAuthConfig{
					Username: getProviderOption(d, "basic_username", "RODE_BASIC_USERNAME"),
					Password: getProviderOption(d, "basic_password", "RODE_BASIC_PASSWORD"),
				},
			}

			rode, err := common.NewRodeClient(
				config,
				grpc.WithUserAgent(provider.UserAgent("terraform-provider-rode", version)),
			)

			return rode, diag.FromErr(err)
		}

		return provider
	}
}

func getProviderOption(d *schema.ResourceData, key, env string) string {
	envVal := os.Getenv(env)
	if envVal != "" {
		return envVal
	}

	return d.Get(key).(string)
}

func getProviderOptionBool(d *schema.ResourceData, key, env string) (bool, error) {
	envVal := os.Getenv(env)
	if envVal != "" {
		return strconv.ParseBool(envVal)
	}

	return d.Get(key).(bool), nil
}
