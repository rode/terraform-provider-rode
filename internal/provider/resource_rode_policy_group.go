package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/proto/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
				Description: "Unique identifier for the policy group",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				// TODO: validate func/regex
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
				Description: "",
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
		if status.Code(err) == codes.NotFound {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.Set("name", policyGroup.Name)
	d.Set("description", policyGroup.Description)
	d.Set("created", formatProtoTimestamp(policyGroup.Created))
	d.Set("updated", formatProtoTimestamp(policyGroup.Updated))

	return nil
}

func resourcePolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("unimplemented")
}

func resourcePolicyGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

	_, err := rode.DeletePolicyGroup(ctx, &v1alpha1.DeletePolicyGroupRequest{Name: d.Id()})
	if status.Code(err) == codes.NotFound {
		d.SetId("")
		return nil
	}

	return diag.FromErr(err)
}
