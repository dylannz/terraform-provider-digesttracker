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

		CreateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			digest := d.Get("digest").(string)
			d.SetId(digest)
			d.Set("version", 1)
			return nil
		},

		ReadContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			return nil
		},

		UpdateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			currentDigest := d.Get("digest").(string)
			previousDigest := d.Id()
			previousVersion := d.Get("version").(int)

			version := previousVersion
			if currentDigest != previousDigest {
				version = previousVersion + 1
			}

			d.SetId(currentDigest)
			d.Set("version", version)
			return nil
		},

		DeleteContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
			d.SetId("")
			return nil
		},
	}
}
