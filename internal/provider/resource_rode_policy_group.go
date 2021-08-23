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
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rode/rode/proto/v1alpha1"
)

var (
	policyNameRegexp                = regexp.MustCompile("^[a-z0-9-_]+$")
	policyNameMessage               = "policy group names may only contain lowercase alphanumeric strings, hyphens, or underscores"
	policyGroupNameValidateDiagFunc = validation.ToDiagFunc(validation.StringMatch(policyNameRegexp, policyNameMessage))
)

func resourcePolicyGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "A policy group is a collection of policies that Rode evaluates against a resource.",
		CreateContext: resourcePolicyGroupCreate,
		ReadContext:   resourcePolicyGroupRead,
		UpdateContext: resourcePolicyGroupUpdate,
		DeleteContext: resourcePolicyGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePolicyGroupImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description:      "Unique identifier for the policy group",
				Type:             schema.TypeString,
				ForceNew:         true,
				Required:         true,
				ValidateDiagFunc: policyGroupNameValidateDiagFunc,
			},
			"description": {
				Description: "A brief summary of the intended use of the policy group",
				Type:        schema.TypeString,
				Optional:    true,
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
				Description: "Indicates that the policy group has been deleted.",
				Computed:    true,
				Type:        schema.TypeBool,
			},
		},
	}
}

func resourcePolicyGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	policyGroup := &v1alpha1.PolicyGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	log.Printf("[DEBUG] Calling CreatePolicyGroup RPC with: %v\n", policyGroup)
	response, err := rode.CreatePolicyGroup(ctx, policyGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Successfully created policy group: %v\n", response)
	d.SetId(response.Name)

	return resourcePolicyGroupRead(ctx, d, meta)
}

func resourcePolicyGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Calling GetPolicyGroup RPC")
	policyGroup, err := rode.GetPolicyGroup(ctx, &v1alpha1.GetPolicyGroupRequest{Name: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", policyGroup.Name)
	d.Set("description", policyGroup.Description)
	d.Set("created", formatProtoTimestamp(policyGroup.Created))
	d.Set("updated", formatProtoTimestamp(policyGroup.Updated))
	d.Set("deleted", policyGroup.Deleted)

	return nil
}

func resourcePolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	policyGroup := &v1alpha1.PolicyGroup{
		Name:        d.Id(),
		Description: d.Get("description").(string),
	}
	log.Printf("[DEBUG] Calling UpdatePolicyGroup RPC with: %v\n", policyGroup)
	response, err := rode.UpdatePolicyGroup(ctx, policyGroup)
	log.Printf("[DEBUG] Successfully updated policy group: %v\n", response)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePolicyGroupRead(ctx, d, meta)
}

func resourcePolicyGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(*rodeClient)
	if err := rode.init(); err != nil {
		return diag.FromErr(err)
	}

	log.Println("[DEBUG] Calling DeletePolicyGroup RPC")
	_, err := rode.DeletePolicyGroup(ctx, &v1alpha1.DeletePolicyGroupRequest{Name: d.Id()})

	return diag.FromErr(err)
}

func resourcePolicyGroupImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	name := d.Id()

	if !policyNameRegexp.MatchString(name) {
		return nil, fmt.Errorf("%s does not match naming restrictions: %s", name, policyNameMessage)
	}

	return []*schema.ResourceData{d}, nil
}
