package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/rode/rode/proto/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
	"testing"
)

func TestAccPolicyGroup_basic(t *testing.T) {
	policyGroup := &v1alpha1.PolicyGroup{
		Name: fmt.Sprintf("tf-acc-%s", strings.ToLower(fake.LetterN(10))),
		Description: fake.LetterN(10),
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProvidersFactory,
		CheckDestroy:      testAccCheckPolicyGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyGroupConfig(policyGroup),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("rode_policy_group.test", "name", policyGroup.Name),
					resource.TestCheckResourceAttr("rode_policy_group.test", "description", policyGroup.Description),
					resource.TestCheckResourceAttrSet("rode_policy_group.test", "created"),
					resource.TestCheckResourceAttrSet("rode_policy_group.test", "updated"),
					testAccCheckPolicyGroupExists("rode_policy_group.test", policyGroup),
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

func testAccCheckPolicyGroupExists(resourceName string, expected *v1alpha1.PolicyGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("policy group not found in state: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("id is not set in state")
		}

		policyGroup, err := rodeClient.GetPolicyGroup(context.Background(), &v1alpha1.GetPolicyGroupRequest{
			Name: rs.Primary.ID,
		})

		if err != nil {
			return err
		}

		if policyGroup.Name != expected.Name {
			return fmt.Errorf("expected policy group name to be %s, got %s", expected.Name, policyGroup.Name)
		}

		if policyGroup.Description != expected.Description {
			return fmt.Errorf("expected policy group description to be %s, got %s", expected.Description, policyGroup.Description)
		}

		if policyGroup.Deleted {
			return fmt.Errorf("policy group was deleted")
		}


		return nil
	}
}

func testAccCheckPolicyGroupDestroyed(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rode_policy_group" {
			continue
		}

		policyGroup, err := rodeClient.GetPolicyGroup(context.Background(), &v1alpha1.GetPolicyGroupRequest{
			Name: rs.Primary.ID,
		})

		if err == nil {
			if policyGroup != nil && policyGroup.Name == rs.Primary.ID && !policyGroup.Deleted {
				return fmt.Errorf("policy group %s still exists", rs.Primary.ID)

			}

			return nil
		}

		if status.Code(err) != codes.NotFound {
			return err
		}
	}
	return nil
}
