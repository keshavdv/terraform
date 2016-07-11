package cloudflare

import (
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudFlareCustomSSL() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFlareCustomSSLCreate,
		Read:   resourceCloudFlareCustomSSLRead,
		Update: resourceCloudFlareCustomSSLUpdate,
		Delete: resourceCloudFlareCustomSSLDelete,

		Schema: map[string]*schema.Schema{
			"certificate": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"private_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"bundle_method": &schema.Schema{
				Optional: true,
				Type:     schema.TypeString,
			},

			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"hosts": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"issuer": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"signature": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"uploaded_on": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"modified_on": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"expires_on": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCloudFlareCustomSSL(d *schema.ResourceData) cloudflare.ZoneCustomSSLOptions {

	newCustomSSL := cloudflare.ZoneCustomSSLOptions{
		Certificate:  d.Get("certificate").(string),
		PrivateKey:   d.Get("private_key").(string),
		BundleMethod: d.Get("bundle_method").(string),
	}

	return newCustomSSL
}

func resourceCloudFlareCustomSSLCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	newCustomSSL := buildCloudFlareCustomSSL(d)
	zoneName := d.Get("domain").(string)

	zoneId, err := client.ZoneIDByName(zoneName)
	if err != nil {
		return fmt.Errorf("Error finding zone %q: %s", zoneName, err)
	}

	d.Set("zone_id", zoneId)

	log.Printf("[DEBUG] CloudFlare CustomSSL create configuration: %#v", newCustomSSL)

	r, err := client.CreateSSL(zoneId, newCustomSSL)
	if err != nil {
		return fmt.Errorf("Failed to create custom SSL: %s", err)
	}

	d.SetId(r.ID)
	log.Printf("[DEBUG] CloudFlare CustomSSL ID: %#v", r.ID)

	return resourceCloudFlareCustomSSLRead(d, meta)
}

func resourceCloudFlareCustomSSLRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	record, err := client.SSLDetails(d.Get("zone_id").(string), d.Id())
	if err != nil {
		return err
	}

	d.SetId(record.ID)
	d.Set("hosts", record.Hosts)
	d.Set("issuer", record.Issuer)
	d.Set("signature", record.Signature)
	d.Set("status", record.Status)
	d.Set("uploaded_on", record.UploadedOn)
	d.Set("modified_on", record.ModifiedOn)
	d.Set("expires_on", record.ExpiresOn)
	d.Set("priority", record.Priority)

	return nil
}

func resourceCloudFlareCustomSSLUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	var err error
	var createdSSL cloudflare.ZoneCustomSSL

	newCustomSSL := buildCloudFlareCustomSSL(d)

	// If only the bundle_method has changed, we need to delete and recreate the SSL certificate
	if d.HasChange("bundle_method") && !d.HasChange("certificate") && !d.HasChange("private_key") {

		log.Printf("[DEBUG] Recreating CloudFlare Custom SSL certificate %s for domain %s", d.Id(), d.Get("domain"))

		err = client.DeleteSSL(d.Get("zone_id").(string), d.Id())
		if err != nil {
			return fmt.Errorf("Error deleting CloudFlare Custom SSL certificate: %s", err)
		}

		createdSSL, err = client.CreateSSL(d.Get("zone_id").(string), newCustomSSL)
		if err != nil {
			return fmt.Errorf("Failed to create custom SSL: %s", err)
		}

	} else {
		createdSSL, err = client.UpdateSSL(d.Get("zone_id").(string), d.Id(), newCustomSSL)
		log.Printf("[DEBUG] Updating CloudFlare Custom SSL certificate %s for domain %s", d.Id(), d.Get("domain"))
		if err != nil {
			return err
		}
	}

	d.SetId(createdSSL.ID)

	return resourceCloudFlareCustomSSLRead(d, meta)
}

func resourceCloudFlareCustomSSLDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	domain := d.Get("domain").(string)

	log.Printf("[INFO] Deleting CloudFlare Custom SSL certificate: %s, %s", domain, d.Id())

	err := client.DeleteSSL(d.Get("zone_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting CloudFlare Custom SSL certificate: %s", err)
	}

	return nil
}
