package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rode/rode/proto/v1alpha1"
	"regexp"
)

var (
	policyNameRegexp                = regexp.MustCompile("^[a-z0-9-_]+$")
	policyGroupNameValidateDiagFunc = validation.ToDiagFunc(validation.StringMatch(policyNameRegexp, "policy group names may only contain lowercase alphanumeric strings, hyphens, or underscores."))
)

func resourcePolicyGroup() *schema.Resource {
	return &schema.Resource{
		Description:   "A policy group is a collection of policies that Rode evaluates against a resource.",
		CreateContext: resourcePolicyGroupCreate,
		ReadContext:   resourcePolicyGroupRead,
		UpdateContext: resourcePolicyGroupUpdate,
		DeleteContext: resourcePolicyGroupDelete,
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
	rode := meta.(v1alpha1.RodeClient)

	policyGroup := &v1alpha1.PolicyGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
	response, err := rode.CreatePolicyGroup(ctx, policyGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.Name)

	return resourcePolicyGroupRead(ctx, d, meta)
}

func resourcePolicyGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

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
	return diag.Errorf("unimplemented")
}

func resourcePolicyGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

	_, err := rode.DeletePolicyGroup(ctx, &v1alpha1.DeletePolicyGroupRequest{Name: d.Id()})

	return diag.FromErr(err)
}
