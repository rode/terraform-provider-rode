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

func TestAccPolicyGroup_basic(t *testing.T) {
	policyGroup := &v1alpha1.PolicyGroup{
		Name:        fmt.Sprintf("tf-acc-%s", strings.ToLower(fake.LetterN(10))),
		Description: fake.LetterN(10),
	}
	var actualPolicyGroup v1alpha1.PolicyGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactory,
		CheckDestroy:      testAccCheckPolicyGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyGroupConfig(policyGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rode_policy_group.test", "name", policyGroup.Name),
					resource.TestCheckResourceAttr("rode_policy_group.test", "description", policyGroup.Description),
					resource.TestCheckResourceAttrSet("rode_policy_group.test", "created"),
					resource.TestCheckResourceAttrSet("rode_policy_group.test", "updated"),
					resource.TestCheckResourceAttr("rode_policy_group.test", "deleted", "false"),
					testAccCheckPolicyGroupExists("rode_policy_group.test", &actualPolicyGroup),
					testAccCheckPolicyGroupValues(policyGroup, &actualPolicyGroup),
				),
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

func testAccCheckPolicyGroupExists(resourceName string, policyGroup *v1alpha1.PolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("policy group not found in state: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set in state")
		}

		rodeClient := testAccProvider.Meta().(v1alpha1.RodeClient)
		group, err := rodeClient.GetPolicyGroup(context.Background(), &v1alpha1.GetPolicyGroupRequest{
			Name: rs.Primary.ID,
		})

		if err != nil {
			return err
		}

		*policyGroup = *group

		return nil
	}
}


func testAccCheckPolicyGroupValues(expected *v1alpha1.PolicyGroup, actual *v1alpha1.PolicyGroup) resource.TestCheckFunc {
	return func(state *terraform.State) error {
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

		rodeClient := testAccProvider.Meta().(v1alpha1.RodeClient)
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
