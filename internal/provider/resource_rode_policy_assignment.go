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
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/proto/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strings"
)

var (
	policyVersionIdValidateDiagFunc = func(v interface{}, p cty.Path) diag.Diagnostics {
		_, err := parsePolicyVersionId(v.(string))

		return diag.FromErr(err)
	}
)

func resourcePolicyAssignment() *schema.Resource {
	return &schema.Resource{
		Description:   "A policy assignment is a mapping between a policy group and a policy version",
		CreateContext: resourcePolicyAssignmentCreate,
		ReadContext:   resourcePolicyAssignmentRead,
		UpdateContext: resourcePolicyAssignmentUpdate,
		DeleteContext: resourcePolicyAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyAssignmentImport,
		},
		Schema: map[string]*schema.Schema{
			"policy_version_id": {
				Description:      "Unique identifier of the versioned policy",
				Required:         true,
				Type:             schema.TypeString,
				ValidateDiagFunc: policyVersionIdValidateDiagFunc,
			},
			"policy_group": {
				Description:      "Name of the policy group to associate with the policy",
				Required:         true,
				ForceNew:         true,
				Type:             schema.TypeString,
				ValidateDiagFunc: policyGroupNameValidateDiagFunc,
			},
			"created": {
				Description: "Creation timestamp",
				Computed:    true,
				Type:        schema.TypeString,
			},
			"updated": {
				Description: "Last updated timestamp",
				Computed:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func resourcePolicyAssignmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}
	policyAssignment := &v1alpha1.PolicyAssignment{
		PolicyVersionId: d.Get("policy_version_id").(string),
		PolicyGroup:     d.Get("policy_group").(string),
	}

	log.Printf("[DEBUG] Calling CreatePolicyAssignment RPC with: %v\n", policyAssignment)
	response, err := rode.CreatePolicyAssignment(ctx, policyAssignment)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Successfully created policy assignment: %v\n", response)
	d.SetId(response.Id)

	return resourcePolicyAssignmentRead(ctx, d, meta)
}

func resourcePolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Calling GetPolicyAssignment RPC")
	response, err := rode.GetPolicyAssignment(ctx, &v1alpha1.GetPolicyAssignmentRequest{Id: d.Id()})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Println("[DEBUG] Policy assignment appears to have been deleted")
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.Set("policy_version_id", response.PolicyVersionId)
	d.Set("policy_group", response.PolicyGroup)
	d.Set("created", formatProtoTimestamp(response.Created))
	d.Set("updated", formatProtoTimestamp(response.Updated))

	return nil
}

func resourcePolicyAssignmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	assignment := &v1alpha1.PolicyAssignment{
		Id:              d.Id(),
		PolicyVersionId: d.Get("policy_version_id").(string),
		PolicyGroup:     d.Get("policy_group").(string),
	}
	log.Printf("[DEBUG] Calling UpdatePolicyAssignment RPC with: %v\n", assignment)
	response, err := rode.UpdatePolicyAssignment(ctx, assignment)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Successfully updated policy assignment: %v\n", response)

	return resourcePolicyAssignmentRead(ctx, d, meta)
}

func resourcePolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Calling DeletePolicyAssignment RPC")
	_, err := rode.DeletePolicyAssignment(ctx, &v1alpha1.DeletePolicyAssignmentRequest{Id: d.Id()})
	if status.Code(err) == codes.NotFound {
		log.Println("[DEBUG] Assignment was already deleted")
		d.SetId("")
		return nil
	}

	return diag.FromErr(err)
}

func resourcePolicyAssignmentImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	assignmentId := d.Id()
	validationMessage := "policy assignment ids should be of the form: policies/$policyId/assignments/$policyGroupName"
	parts := strings.Split(assignmentId, "/")
	if len(parts) != 4 {
		return nil, fmt.Errorf(validationMessage)
	}

	if parts[0] != "policies" || parts[2] != "assignments" {
		return nil, fmt.Errorf(validationMessage)
	}

	policyId := parts[1]
	if _, err := uuid.ParseUUID(policyId); err != nil {
		return nil, fmt.Errorf("invalid policy id: %s", err)
	}

	policyGroupName := parts[3]
	if !policyNameRegexp.MatchString(policyGroupName) {
		return nil, fmt.Errorf("policy group name '%s' does not match naming restrictions", policyGroupName)
	}

	return []*schema.ResourceData{d}, nil
}
