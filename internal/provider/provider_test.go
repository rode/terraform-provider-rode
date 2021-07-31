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
