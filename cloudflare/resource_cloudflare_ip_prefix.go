package cloudflare

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/pkg/errors"
)

func resourceCloudflareIPPrefix() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudflareIPPrefixCreate,
		Read:   resourceCloudflareIPPrefixRead,
		Update: resourceCloudflareIPPrefixUpdate,
		Delete: resourceCloudflareIPPrefixDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloudflareIPPrefixImport,
		},

		Schema: map[string]*schema.Schema{
			"prefix_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"advertisement": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on", "off"}, false),
				Computed:     true,
				Optional:     true,
			},
		},
	}
}

func resourceCloudflareIPPrefixCreate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceCloudflareIPPrefixRead(d, meta); err != nil {
		return err
	}

	return resourceCloudflareIPPrefixUpdate(d, meta)
}

func resourceCloudflareIPPrefixImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	prefixID := d.Id()
	d.Set("prefix_id", prefixID)

	resourceCloudflareIPPrefixRead(d, meta)

	return []*schema.ResourceData{d}, nil
}

func resourceCloudflareIPPrefixRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	prefixID := d.Get("prefix_id").(string)
	d.SetId(prefixID)

	prefix, err := client.GetPrefix(context.Background(), d.Id())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading IP prefix information for %q", d.Id()))
	}

	d.Set("description", prefix.Description)

	advertisementStatus, err := client.GetAdvertisementStatus(context.Background(), d.Id())
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Error reading advertisement status of IP prefix for %q", d.Id()))
	}

	d.Set("advertisement", stringFromBool(advertisementStatus.Advertised))

	return nil
}

func resourceCloudflareIPPrefixUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	prefixID := d.Get("prefix_id").(string)
	d.SetId(prefixID)

	if _, ok := d.GetOk("description"); ok && d.HasChange("description") {
		if _, err := client.UpdatePrefixDescription(context.Background(), d.Id(), d.Get("description").(string)); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Cannot update prefix description for %q", d.Id()))
		}
	}

	if _, ok := d.GetOk("advertisement"); ok && d.HasChange("advertisement") {
		os.Exit(1)
		if _, err := client.UpdateAdvertisementStatus(context.Background(), d.Id(), boolFromString(d.Get("advertisement").(string))); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Cannot update prefix advertisement status for %q", d.Id()))
		}
	}

	return nil
}

// Deletion of prefixes is not really supported, so we keep this as a dummy
func resourceCloudflareIPPrefixDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
