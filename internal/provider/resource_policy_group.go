package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/proto/v1alpha1"
	"strconv"
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

func policyGroupToData(policyGroup *v1alpha1.PolicyGroup, d *schema.ResourceData) {
	d.SetId(policyGroup.Name)

	d.Set("name", policyGroup.Name)
	d.Set("description", policyGroup.Description)
	d.Set("created", formatProtoTimestamp(policyGroup.Created))
	d.Set("updated", formatProtoTimestamp(policyGroup.Updated))
	d.Set("deleted", strconv.FormatBool(policyGroup.Deleted))
}

func dataToPolicyGroup(d *schema.ResourceData, policyGroup *v1alpha1.PolicyGroup) error {
	policyGroup.Name = d.Get("name").(string)
	policyGroup.Description = d.Get("description").(string)

	if value, ok := d.GetOk("created"); ok {
		created, err := asProtoTimestamp(value)
		if err != nil {
			return fmt.Errorf(`error parsing "created" timestamp: %s`, err)
		}
		policyGroup.Created = created
	}

	if value, ok := d.GetOk("updated"); ok {
		updated, err := asProtoTimestamp(value)
		if err != nil {
			return fmt.Errorf(`error parsing "updated" timestamp: %s`, err)
		}
		policyGroup.Updated = updated
	}

	if value, ok := d.GetOk("deleted"); ok {
		deleted, err := strconv.ParseBool(value.(string))
		if err != nil {
			return fmt.Errorf(`error parsing value of "deleted" field: %s`, err)
		}
		policyGroup.Deleted = deleted
	}

	return nil
}

func resourcePolicyGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

	policyGroup := &v1alpha1.PolicyGroup{}
	if err := dataToPolicyGroup(d, policyGroup); err != nil {
		return diag.FromErr(err)
	}

	response, err := rode.CreatePolicyGroup(ctx, policyGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	policyGroupToData(response, d)

	return nil
}

func resourcePolicyGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

	policyGroup, err := rode.GetPolicyGroup(ctx, &v1alpha1.GetPolicyGroupRequest{Name: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	policyGroupToData(policyGroup, d)

	return nil
}

func resourcePolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("unimplemented")
}

func resourcePolicyGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

	_, err := rode.DeletePolicyGroup(ctx, &v1alpha1.DeletePolicyGroupRequest{Name: d.Id()})

	if isDeletionError(err) {
		return diag.FromErr(err)
	}

	return nil
}
