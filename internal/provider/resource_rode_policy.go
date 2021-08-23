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
	"log"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/proto/v1alpha1"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "An Open Policy Agent Rego policy.",
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyImport,
		},
		CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, i interface{}) error {
			if diff.HasChange("rego_content") {
				return diff.SetNewComputed("policy_version_id")
			}

			return nil
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Policy name",
				Required:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Description: "A brief summary of the policy",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"current_version": {
				Description: "Current version of the policy",
				Computed:    true,
				Type:        schema.TypeInt,
			},
			"policy_version_id": {
				Computed:    true,
				Description: "Policy version id",
				Type:        schema.TypeString,
			},
			"message": {
				Description: "A summary of changes since the last version",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"rego_content": {
				Description: "The Rego code",
				Type:        schema.TypeString,
				Required:    true,
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
			"deleted": {
				Description: "Indicates that the policy has been deleted.",
				Computed:    true,
				Type:        schema.TypeBool,
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	policy := &v1alpha1.Policy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Policy: &v1alpha1.PolicyEntity{
			Message:     d.Get("message").(string),
			RegoContent: d.Get("rego_content").(string),
		},
	}

	log.Printf("[DEBUG] Calling CreatePolicy RPC with: %v\n", policy)
	response, err := rode.CreatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Successfully created policy: %v\n", response)

	d.SetId(response.Id)

	return resourcePolicyRead(ctx, d, meta)
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Calling GetPolicy RPC")
	policy, err := rode.GetPolicy(ctx, &v1alpha1.GetPolicyRequest{Id: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("current_version", policy.CurrentVersion)
	d.Set("policy_version_id", policy.Policy.Id)
	d.Set("message", policy.Policy.Message)
	d.Set("rego_content", policy.Policy.RegoContent)
	d.Set("created", formatProtoTimestamp(policy.Created))
	d.Set("updated", formatProtoTimestamp(policy.Updated))
	d.Set("deleted", policy.Deleted)

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	policy := &v1alpha1.Policy{
		Id:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Policy: &v1alpha1.PolicyEntity{
			Message:     d.Get("message").(string),
			RegoContent: d.Get("rego_content").(string),
		},
	}

	log.Printf("[DEBUG] Calling UpdatePolicy RPC with: %v\n", policy)
	response, err := rode.UpdatePolicy(ctx, &v1alpha1.UpdatePolicyRequest{
		Policy: policy,
	})

	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Successfully updated policy: %v\n", response)

	return resourcePolicyRead(ctx, d, meta)
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Calling DeletePolicy RPC")
	_, err := rode.DeletePolicy(ctx, &v1alpha1.DeletePolicyRequest{
		Id: d.Id(),
	})

	return diag.FromErr(err)
}

func resourcePolicyImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	policyId := d.Id()
	if _, err := uuid.ParseUUID(policyId); err != nil {
		return nil, fmt.Errorf("invalid policy id: %s", err)
	}

	return []*schema.ResourceData{d}, nil
}
