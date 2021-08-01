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
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rode/rode/proto/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	"strings"
)

var (
	isUuidValidateDiagFunc          = validation.ToDiagFunc(validation.IsUUID)
	policyVersionIdValidateDiagFunc = func(v interface{}, p cty.Path) diag.Diagnostics {
		policyVersionId := v.(string)

		parts := strings.Split(policyVersionId, ".")
		if len(parts) != 2 {
			return diag.Errorf("policy version id does not match format")
		}

		policyId := parts[0]
		_, err := strconv.Atoi(parts[1])
		if err != nil {
			return diag.Errorf("policy version id does not contain a version: %s", err)
		}

		return isUuidValidateDiagFunc(policyId, p)
	}
)

func resourcePolicyAssignment() *schema.Resource {
	return &schema.Resource{
		Description:   "A policy assignment is a mapping between a policy group and a policy version",
		CreateContext: resourcePolicyAssignmentCreate,
		ReadContext:   resourcePolicyAssignmentRead,
		UpdateContext: resourcePolicyAssignmentUpdate,
		DeleteContext: resourcePolicyAssignmentDelete,
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

	response, err := rode.CreatePolicyAssignment(ctx, policyAssignment)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.Id)

	return resourcePolicyAssignmentRead(ctx, d, meta)
}

func resourcePolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	response, err := rode.GetPolicyAssignment(ctx, &v1alpha1.GetPolicyAssignmentRequest{Id: d.Id()})

	if err != nil {
		if status.Code(err) == codes.NotFound {
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
	return diag.Errorf("unimplemented")
}

func resourcePolicyAssignmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	_, err := rode.DeletePolicyAssignment(ctx, &v1alpha1.DeletePolicyAssignmentRequest{Id: d.Id()})
	if status.Code(err) == codes.NotFound {
		d.SetId("")
		return nil
	}

	return diag.FromErr(err)
}
