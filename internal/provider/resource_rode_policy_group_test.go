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
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rode/rode/proto/v1alpha1"
)

func TestAccPolicyGroup_basic(t *testing.T) {
	resourceName := "rode_policy_group.test"
	policyGroup := &v1alpha1.PolicyGroup{
		Name:        fmt.Sprintf("tf-acc-%s", strings.ToLower(fake.LetterN(10))),
		Description: fake.LetterN(10),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactory,
		CheckDestroy:      testAccCheckPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyGroupConfig(policyGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policyGroup.Name),
					resource.TestCheckResourceAttr(resourceName, "description", policyGroup.Description),
					resource.TestCheckResourceAttrSet(resourceName, "created"),
					resource.TestCheckResourceAttrSet(resourceName, "updated"),
					resource.TestCheckResourceAttr(resourceName, "deleted", "false"),
					testAccCheckPolicyGroupExists(resourceName, policyGroup),
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

func TestAccPolicyGroup_update(t *testing.T) {
	resourceName := "rode_policy_group.test"
	policyGroupName := fmt.Sprintf("tf-acc-%s", strings.ToLower(fake.LetterN(10)))
	policyGroup := &v1alpha1.PolicyGroup{
		Name:        policyGroupName,
		Description: fake.LetterN(10),
	}
	updatedPolicyGroup := &v1alpha1.PolicyGroup{
		Name:        policyGroupName,
		Description: fake.LetterN(10),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactory,
		CheckDestroy:      testAccCheckPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyGroupConfig(policyGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policyGroup.Name),
					resource.TestCheckResourceAttr(resourceName, "description", policyGroup.Description),
					testAccCheckPolicyGroupExists(resourceName, policyGroup),
				),
			},
			{
				Config: testAccPolicyGroupConfig(updatedPolicyGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updatedPolicyGroup.Name),
					resource.TestCheckResourceAttr(resourceName, "description", updatedPolicyGroup.Description),
					testAccCheckPolicyGroupExists(resourceName, updatedPolicyGroup),
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

func TestAccPolicyGroup_validate_name(t *testing.T) {
	policyGroup := &v1alpha1.PolicyGroup{
		Name: fmt.Sprintf("tf-acc-%s!@#$", strings.ToUpper(fake.LetterN(10))),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactory,
		Steps: []resource.TestStep{
			{
				Config:      testAccPolicyGroupConfig(policyGroup),
				ExpectError: regexp.MustCompile("policy group names may only contain"),
			},
		},
	})
}

func testAccPolicyGroupConfig(policyGroup *v1alpha1.PolicyGroup) string {
	return fmt.Sprintf(`
resource "rode_policy_group" "test" {
  name        = "%s"
  description = "%s"
}
`, policyGroup.Name, policyGroup.Description)
}

func testAccCheckPolicyGroupExists(resourceName string, expected *v1alpha1.PolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("policy group not found in state: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set in state")
		}

		rodeClient := testAccProvider.Meta().(*rodeClient)
		actual, err := rodeClient.GetPolicyGroup(context.Background(), &v1alpha1.GetPolicyGroupRequest{
			Name: rs.Primary.ID,
		})

		if err != nil {
			return err
		}

		if expected.Name != actual.Name {
			return fmt.Errorf("expected policy group name to be '%s', but was '%s'", expected.Name, actual.Name)
		}

		if expected.Description != actual.Description {
			return fmt.Errorf("expected policy group description to be '%s', but was '%s'", expected.Description, actual.Description)
		}

		return nil
	}
}

func testAccCheckPolicyGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rode_policy_group" {
			continue
		}

		rodeClient := testAccProvider.Meta().(*rodeClient)
		policyGroup, err := rodeClient.GetPolicyGroup(context.Background(), &v1alpha1.GetPolicyGroupRequest{
			Name: rs.Primary.ID,
		})

		if err != nil {
			return err
		}

		if policyGroup != nil && !policyGroup.Deleted {
			return fmt.Errorf("policy group %s still exists", rs.Primary.ID)
		}
	}
	return nil
}
