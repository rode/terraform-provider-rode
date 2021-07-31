provider "rode" {
  // RODE_HOST
  host = "localhost:50051"
  // RODE_DISABLE_TRANSPORT_SECURITY
  disable_transport_security = true

  // basic and oidc configuration is optional
  // only one authentication method can be configured

  // RODE_OIDC_CLIENT_ID
  oidc_client_id = "terraform"
  // RODE_OIDC_CLIENT_SECRET
  oidc_client_secret = "top secret"
  // RODE_OIDC_TOKEN_URL
  oidc_token_url = "https://idp.example.com/oauth2/token"
  // RODE_OIDC_SCOPES
  oidc_scopes = "rode terraform"
  // RODE_OIDC_TLS_INSECURE_SKIP_VERIFY
  oidc_tls_insecure_skip_verify = false
  // RODE_BASIC_USERNAME
  basic_username = "policy-administrator"
  // RODE_BASIC_PASSWORD
  basic_password = "password"
}
