package jamf

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/yohan460/go-jamf-api"
)

func resourceJamfSite() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJamfSiteCreate,
		ReadContext:   resourceJamfSiteRead,
		UpdateContext: resourceJamfSiteUpdate,
		DeleteContext: resourceJamfSiteDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Read:   schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: importJamfSiteState,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildJamfSiteStruct(d *schema.ResourceData) *jamf.Site {
	var out jamf.Site
	out.SetId(d.Id())
	out.SetName(d.Get("name").(string))

	return &out
}

func resourceJamfSiteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*jamf.Client)

	b := buildJamfSiteStruct(d)

	out, err := c.CreateSite(b.Name, b.StreetAddress1, b.StreetAddress2, b.City, b.StateProvince, b.ZipPostalCode, b.Country)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(out.GetId())

	return resourceJamfSiteRead(ctx, d, m)
}

func resourceJamfSiteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*jamf.Client)

	resp, err := c.GetSiteByName(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", resp.GetName())

	return diags
}

func resourceJamfSiteUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*jamf.Client)

	b := buildJamfSiteStruct(d)
	d.SetId(b.GetId())

	if _, err := c.UpdateSite(b); err != nil {
		return diag.FromErr(err)
	}

	return resourceJamfSiteRead(ctx, d, m)
}

func resourceJamfSiteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*jamf.Client)
	b := buildJamfSiteStruct(d)

	if err := c.DeleteSite(*b.Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}

func importJamfSiteState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	c := m.(*jamf.Client)
	d.SetId(d.Id())
	resp, err := c.GetSite(d.Id())
	if err != nil {
		return nil, fmt.Errorf("cannot get Site data")
	}

	d.Set("name", resp.GetName())

	return []*schema.ResourceData{d}, nil
}
