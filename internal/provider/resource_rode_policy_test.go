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
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rode/rode/proto/v1alpha1"
	"strings"
	"testing"
)

func TestAccPolicy_basic(t *testing.T) {
	resourceName := "rode_policy.test"
	policy := &v1alpha1.Policy{
		Name:        fmt.Sprintf("tf-acc-%s", fake.LetterN(10)),
		Description: fake.LetterN(10),
		Policy: &v1alpha1.PolicyEntity{
			RegoContent: minimalPolicy,
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactory,
		CheckDestroy:      testAccPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig(policy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policy.Name),
					resource.TestCheckResourceAttr(resourceName, "description", policy.Description),
					resource.TestCheckResourceAttr(resourceName, "current_version", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_version_id"),
					resource.TestCheckResourceAttr(resourceName, "message", "Initial policy creation"),
					resource.TestCheckResourceAttrSet(resourceName, "rego_content"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
					resource.TestCheckResourceAttr(resourceName, "deleted", "false"),
					testAccPolicyExists(resourceName, policy),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPolicyConfig(policy *v1alpha1.Policy) string {
	return fmt.Sprintf(`
resource "rode_policy" "test" {
	name  		 = "%s"
    description  = "%s"
	message      = "%s"
	rego_content = <<EOF
%s
EOF
}
`,
		policy.Name,
		policy.Description,
		policy.Policy.Message,
		policy.Policy.RegoContent,
	)
}

func testAccPolicyExists(resourceName string, expected *v1alpha1.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("policy not found in state: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set in state")
		}

		rodeClient := testAccProvider.Meta().(*rodeClient)
		actual, err := rodeClient.GetPolicy(context.Background(), &v1alpha1.GetPolicyRequest{
			Id: rs.Primary.ID,
		})
		if err != nil {
			return err
		}

		if actual.Name != expected.Name {
			return fmt.Errorf("expected policy name to equal '%s' but was '%s'", expected.Name, actual.Name)
		}

		if actual.Description != expected.Description {
			return fmt.Errorf("expected policy description to equal '%s' but was '%s'", expected.Description, actual.Description)
		}

		if actual.CurrentVersion != 1 {
			return fmt.Errorf("policy should be at initial version but was at %d", actual.CurrentVersion)
		}

		if strings.TrimSpace(actual.Policy.RegoContent) != strings.TrimSpace(expected.Policy.RegoContent) {
			return fmt.Errorf("expected policy Rego to equal '%s' but was '%s'", expected.Policy.RegoContent, actual.Policy.RegoContent)
		}

		return nil
	}
}

func testAccPolicyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rode_policy" {
			continue
		}

		rodeClient := testAccProvider.Meta().(*rodeClient)
		policy, err := rodeClient.GetPolicy(context.Background(), &v1alpha1.GetPolicyRequest{
			Id: rs.Primary.ID,
		})

		if err != nil {
			return err
		}

		if policy != nil && !policy.Deleted {
			return fmt.Errorf("policy '%s' still exists", rs.Primary.ID)
		}
	}

	return nil
}
