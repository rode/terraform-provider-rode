package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/proto/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
				Description: "Unique identifier of the versioned policy",
				Required:    true,
				Type:        schema.TypeString,
			},
			"policy_group": {
				Description: "Name of the policy group to associate with the policy",
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
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
	rode := meta.(v1alpha1.RodeClient)
	policyAssignment := &v1alpha1.PolicyAssignment{
		PolicyVersionId: d.Get("policy_version_id").(string),
		PolicyGroup: d.Get("policy_group").(string),
	}

	response, err := rode.CreatePolicyAssignment(ctx, policyAssignment)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.Id)

	return resourcePolicyAssignmentRead(ctx, d, meta)
}

func resourcePolicyAssignmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

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
	rode := meta.(v1alpha1.RodeClient)

	_, err := rode.DeletePolicyAssignment(ctx, &v1alpha1.DeletePolicyAssignmentRequest{Id: d.Id()})
	if status.Code(err) == codes.NotFound {
		d.SetId("")
		return nil
	}

	return diag.FromErr(err)
}
