package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rode/rode/proto/v1alpha1"
	"strconv"
)

func resourcePolicy() *schema.Resource {
	return &schema.Resource{
		Description:   "An Open Policy Agent Rego policy.",
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Unique identifier for the policy",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Policy name",
				Required:    true,
				Type:        schema.TypeString,
			},
			"description": {
				Description: "A brief summary of the policy",
				Required:    true,
				Type:        schema.TypeString,
			},
			"current_version": {
				Description: "Current version of the policy",
				Computed:    true,
				Type:        schema.TypeInt,
			},
			"policy": {
				Type:        schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Description: "The versioned Rego policy",
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Policy version id",
							Computed:    true,
							Type:        schema.TypeString,
						},
						"version": {
							Description: "Numeric version",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"message": {
							Description: "A summary of changes since the last version",
							Type:        schema.TypeString,
							Optional:    true,
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
					},
				},
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

func policyToData(policy *v1alpha1.Policy, d *schema.ResourceData) {
	d.SetId(policy.Id)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("current_version", policy.CurrentVersion)
	d.Set("policy", map[string]interface{}{
		"id":           policy.Policy.Id,
		"version":      policy.Policy.Version,
		"message":      policy.Policy.Message,
		"rego_content": policy.Policy.RegoContent,
		"created":      formatProtoTimestamp(policy.Policy.Created),
	})
	d.Set("created", formatProtoTimestamp(policy.Created))
	d.Set("updated", formatProtoTimestamp(policy.Updated))
	d.Set("deleted", strconv.FormatBool(policy.Deleted))
}

func dataToPolicy(d *schema.ResourceData, policy *v1alpha1.Policy) error {
	policy.Id = d.Id()
	policy.Name = d.Get("name").(string)
	policy.Description = d.Get("description").(string)

	if value, ok := d.GetOk("current_version"); ok {
		if version, ok := value.(uint32); ok {
			policy.CurrentVersion = version
		}
	}

	if value, ok := d.GetOk("policy"); ok {
		policyVersionData := value.([]interface{})[0].(map[string]interface{})

		policyVersion := &v1alpha1.PolicyEntity{}
		policyVersion.Id = policyVersionData["id"].(string)
		policyVersion.Message = policyVersionData["message"].(string)
		policyVersion.RegoContent = policyVersionData["rego_content"].(string)

		// TODO: version and created are unset
		if v, ok := policyVersionData["version"]; ok {
			if version, ok := v.(uint32); ok {
				policyVersion.Version = version
			}
		}

		if v, ok := policyVersionData["created"]; ok && v != "" {
			created, err := asProtoTimestamp(v)
			if err != nil {
				return err
			}

			policyVersion.Created = created
		}

		policy.Policy = policyVersion
	}

	if value, ok := d.GetOk("created"); ok {
		created, err := asProtoTimestamp(value)
		if err != nil {
			return fmt.Errorf(`error parsing "created" timestamp: %s`, err)
		}
		policy.Created = created
	}

	if value, ok := d.GetOk("updated"); ok {
		updated, err := asProtoTimestamp(value)
		if err != nil {
			return fmt.Errorf(`error parsing "updated" timestamp: %s`, err)
		}
		policy.Updated = updated
	}

	if value, ok := d.GetOk("deleted"); ok {
		deleted, err := strconv.ParseBool(value.(string))
		if err != nil {
			return fmt.Errorf(`error parsing value of "deleted" field: %s`, err)
		}
		policy.Deleted = deleted
	}

	return nil
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)
	policy := &v1alpha1.Policy{}
	if err := dataToPolicy(d, policy); err != nil {
		return diag.FromErr(err)
	}

	response, err := rode.CreatePolicy(ctx, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	policyToData(response, d)

	return nil
}

func resourcePolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

	policy, err := rode.GetPolicy(ctx, &v1alpha1.GetPolicyRequest{Id: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}

	policyToData(policy, d)

	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("unimplemented")
}

func resourcePolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	rode := meta.(v1alpha1.RodeClient)

	_, err := rode.DeletePolicy(ctx, &v1alpha1.DeletePolicyRequest{
		Id: d.Id(),
	})

	if isDeletionError(err) {
		return diag.FromErr(err)
	}

	return nil
}

