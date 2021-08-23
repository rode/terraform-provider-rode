// Copyright 2021 The Rode Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/common"
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
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("RODE_HOST", ""),
				},
				"disable_transport_security": {
					Description: "Disables transport security for the gRPC connection to Rode. Can also be set with the `RODE_DISABLE_TRANSPORT_SECURITY` environment variable.",
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RODE_DISABLE_TRANSPORT_SECURITY", false),
				},
				"lazy_init": {
					Description: "Defers instantiation of the Rode client until the first time the provider is used. This can be useful when provider config depends on other resources being applied.",
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RODE_LAZY_INIT", false),
				},
				"oidc_client_id": {
					Description: "OIDC/OAuth2 client id that is permitted the client credentials grant. Can be set with the `RODE_OIDC_CLIENT_ID` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RODE_OIDC_CLIENT_ID", ""),
				},
				"oidc_client_secret": {
					Description: "Corresponding client secret for oidc_client_id. Can be set with the `RODE_OIDC_CLIENT_SECRET` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RODE_OIDC_CLIENT_SECRET", ""),
					Sensitive:   true,
				},
				"oidc_token_url": {
					Description: "OAuth2 token url. Can be set with the OIDC_TOKEN_URL environment variable",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RODE_OIDC_TOKEN_URL", ""),
				},
				"oidc_scopes": {
					Description: "A space-delimited list of scopes to request in the client credentials grant. Can also be set with the `RODE_OIDC_SCOPES` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RODE_OIDC_SCOPES", ""),
				},
				"oidc_tls_insecure_skip_verify": {
					Description: "Disable transport security when communicating with the OAuth2 server. Only recommended for local development. Set with the `RODE_OIDC_TLS_INSECURE_SKIP_VERIFY` environment variable.",
					Type:        schema.TypeBool,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("OIDC_TLS_INSECURE_SKIP_VERIFY", false),
				},
				"basic_username": {
					Description: "The username configured in the Rode instance for basic auth. Cannot be configured alongside any of the OIDC options. Can be set with the `RODE_BASIC_USERNAME` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RODE_BASIC_USERNAME", ""),
				},
				"basic_password": {
					Description: "Corresponding password for basic_username. Can be set with the `RODE_BASIC_PASSWORD` environment variable.",
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("RODE_BASIC_PASSWORD", ""),
					Sensitive:   true,
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"rode_policy_group":      resourcePolicyGroup(),
				"rode_policy":            resourcePolicy(),
				"rode_policy_assignment": resourcePolicyAssignment(),
			},
		}

		provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			log.Println("[DEBUG] Provider configure called")
			config := &common.ClientConfig{
				Rode: &common.RodeClientConfig{
					Host:                     d.Get("host").(string),
					DisableTransportSecurity: d.Get("disable_transport_security").(bool),
				},
				OIDCAuth: &common.OIDCAuthConfig{
					ClientID:              d.Get("oidc_client_id").(string),
					ClientSecret:          d.Get("oidc_client_secret").(string),
					TokenURL:              d.Get("oidc_token_url").(string),
					TlsInsecureSkipVerify: d.Get("oidc_tls_insecure_skip_verify").(bool),
					Scopes:                d.Get("oidc_scopes").(string),
				},
				BasicAuth: &common.BasicAuthConfig{
					Username: d.Get("basic_username").(string),
					Password: d.Get("basic_password").(string),
				},
			}

			rodeClient := &rodeClient{
				config:    config,
				userAgent: provider.UserAgent("terraform-provider-rode", version),
			}

			lazyInit := d.Get("lazy_init").(bool)
			if !lazyInit {
				log.Println("[DEBUG] Lazy initialization is disabled, instantiating Rode client immediately")
				if err := rodeClient.init(); err != nil {
					return nil, diag.FromErr(err)
				}
			}

			log.Println("[DEBUG] Delaying Rode client initialization because lazy_init is set")
			return rodeClient, nil
		}

		return provider
	}
}
