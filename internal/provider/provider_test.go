package provider

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"testing"
)

var (
	testAccProvider         *schema.Provider
	testAccProvidersFactory map[string]func() (*schema.Provider, error)
	fake                    *gofakeit.Faker
)

func init() {
	fake = gofakeit.New(0)
	testAccProvider = New("acceptance")()
	testAccProvidersFactory = map[string]func() (*schema.Provider, error){
		"rode": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("RODE_HOST") == "" {
		t.Fatal("RODE_HOST must be set for acceptance tests")
	}
}
