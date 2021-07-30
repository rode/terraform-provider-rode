package provider

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/common"
	"github.com/rode/rode/proto/v1alpha1"
	"os"
	"testing"
)

var (
	testAccProvidersFactory map[string]func() (*schema.Provider, error)
	rodeClient              v1alpha1.RodeClient
	fake                    *gofakeit.Faker
)

func init() {
	fake = gofakeit.New(0)
	testAccProvidersFactory = map[string]func() (*schema.Provider, error){
		"rode": func() (*schema.Provider, error) {
			return New(), nil
		},
	}

	// TODO: find a way to get the provider's client instead
	rodeClient, _ = common.NewRodeClient(&common.ClientConfig{
		Rode: &common.RodeClientConfig{
			Host:                     os.Getenv("RODE_HOST"),
			DisableTransportSecurity: true,
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if rodeClient == nil {
		t.Fatal("Failed to instantiate test client")
	}

	if os.Getenv("RODE_HOST") == "" {
		t.Fatal("RODE_HOST must be set for acceptance tests")
	}
}
