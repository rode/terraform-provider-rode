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
	_ "embed"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rode/rode/proto/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
	"strings"
	"testing"
)

func TestAccPolicyAssignment_basic(t *testing.T) {
	resourceName := "rode_policy_assignment.test"
	policy := &v1alpha1.Policy{
		Name: fmt.Sprintf("tf-acc-%s", fake.LetterN(10)),
		Policy: &v1alpha1.PolicyEntity{
			RegoContent: minimalPolicy,
		},
	}
	policyGroup := &v1alpha1.PolicyGroup{
		Name: fmt.Sprintf("tf-acc-%s", strings.ToLower(fake.LetterN(10))),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactory,
		CheckDestroy:      testAccPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyAssignmentFullConfig(policy, policyGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "policy_group", policyGroup.Name),
					resource.TestCheckResourceAttrSet(resourceName, "policy_version_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
					testAccCheckPolicyAssignmentExists(resourceName, policyGroup.Name),
				),
			},
		},
	})
}

func TestAccPolicyAssignment_invalid_policy_group(t *testing.T) {
	policyVersionId := fmt.Sprintf("%s.%d", fake.UUID(), fake.Number(1, 5))
	policyGroupName := fmt.Sprintf("tf-acc-%s-$!@", strings.ToUpper(fake.LetterN(10)))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactory,
		CheckDestroy:      testAccPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPolicyAssignmentConfig(policyVersionId, policyGroupName),
				ExpectError: regexp.MustCompile("policy group names may only contain"),
			},
		},
	})
}

func TestAccPolicyAssignment_invalid_policy_version_id(t *testing.T) {
	policyVersionId := fake.LetterN(10)
	policyGroupName := fmt.Sprintf("tf-acc-%s", strings.ToLower(fake.LetterN(10)))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactory,
		CheckDestroy:      testAccPolicyAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccPolicyAssignmentConfig(policyVersionId, policyGroupName),
				ExpectError: regexp.MustCompile("policy version id does not match format"),
			},
		},
	})
}

func testAccPolicyAssignmentFullConfig(policy *v1alpha1.Policy, policyGroup *v1alpha1.PolicyGroup) string {
	return fmt.Sprintf(`
resource "rode_policy_group" "test" {
	name 		= "%s"
	description = "%s"
}

resource "rode_policy" "test" {
	name  		 = "%s"
	rego_content = <<EOF
%s
EOF
}

resource "rode_policy_assignment" "test" {
	policy_version_id = rode_policy.test.policy_version_id
	policy_group      = rode_policy_group.test.name
}
`,
		policyGroup.Name,
		policyGroup.Description,
		policy.Name,
		policy.Policy.RegoContent)
}

func testAccPolicyAssignmentConfig(policyVersionId, policyGroup string) string {
	return fmt.Sprintf(`
resource "rode_policy_assignment" "test" {
 	policy_version_id = "%s"
	policy_group      = "%s"
}
`, policyVersionId, policyGroup)
}

func testAccCheckPolicyAssignmentExists(resourceName, policyGroupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("policy assignment not found in state: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set in state")
		}

		rodeClient := testAccProvider.Meta().(v1alpha1.RodeClient)
		actual, err := rodeClient.GetPolicyAssignment(context.Background(), &v1alpha1.GetPolicyAssignmentRequest{
			Id: rs.Primary.ID,
		})
		if err != nil {
			return err
		}

		if actual.PolicyGroup != policyGroupName {
			return fmt.Errorf("expected assignment to group '%s', but was '%s'", policyGroupName, actual.PolicyGroup)
		}

		return nil
	}
}

func testAccPolicyAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rode_policy_assignment" {
			continue
		}

		rodeClient := testAccProvider.Meta().(v1alpha1.RodeClient)
		_, err := rodeClient.GetPolicyAssignment(context.Background(), &v1alpha1.GetPolicyAssignmentRequest{
			Id: rs.Primary.ID,
		})

		if err != nil {
			if status.Code(err) == codes.NotFound {
				return nil
			}

			return err
		}

		return fmt.Errorf("policy assignment %s still exists", rs.Primary.ID)
	}

	return nil
}
