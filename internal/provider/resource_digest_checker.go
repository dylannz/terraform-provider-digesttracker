package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDigestTracker() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"digest": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},

		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
			// Only handle version calculation if digest has changed
			if !d.HasChange("digest") {
				return nil
			}

			// Try to read prior version from state (old), default 0
			prior := 0
			if oldV, _ := d.GetChange("version"); oldV != nil {
				if vInt, ok := oldV.(int); ok {
					prior = vInt
				}
			} else if v, ok := d.GetOk("version"); ok { // fallback
				prior = v.(int)
			}

			return d.SetNew("version", prior+1)
		},

		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			digest := d.Get("digest").(string)
			d.SetId(digest)
			return nil
		},

		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return nil
		},

		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			currentDigest := d.Get("digest").(string)
			d.SetId(currentDigest)
			return nil
		},

		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			d.SetId("")
			return nil
		},
	}
}
