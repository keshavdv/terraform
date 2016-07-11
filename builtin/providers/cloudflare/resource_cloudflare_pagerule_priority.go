package cloudflare

import (
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudFlarePageRulePriority() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFlarePageRulePriorityCreate,
		Read:   resourceCloudFlarePageRulePriorityRead,
		Update: resourceCloudFlarePageRulePriorityUpdate,
		Delete: resourceCloudFlarePageRulePriorityDelete,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"rules": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudFlarePageRulePriorityCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	d.SetId(randomHash())

	zoneName := d.Get("domain").(string)
	zoneId, err := client.ZoneIDByName(zoneName)
	if err != nil {
		return fmt.Errorf("Error finding zone %q: %s", zoneName, err)
	}

	d.Set("zone_id", zoneId)
	return resourceCloudFlarePageRulePriorityUpdate(d, meta)
}

func resourceCloudFlarePageRulePriorityRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloudFlarePageRulePriorityUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	rules := d.Get("rules").([]interface{})
	priority := len(rules)
	for _, ruleId := range rules {
		oldRule, err := client.PageRule(d.Get("zone_id").(string), ruleId.(string))
		if err != nil {
			return err
		}

		oldRule.Priority = cloudflare.MaybeInt(priority)
		_, err = client.ChangePageRule(d.Get("zone_id").(string), oldRule.ID, oldRule)
		if err != nil {
			return fmt.Errorf("Error updating PageRule priority: %s", err)
		}
		priority--
	}

	return nil
}

func resourceCloudFlarePageRulePriorityDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
