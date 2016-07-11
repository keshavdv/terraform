package cloudflare

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudFlarePageRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFlarePageRuleCreate,
		Read:   resourceCloudFlarePageRuleRead,
		Update: resourceCloudFlarePageRuleUpdate,
		Delete: resourceCloudFlarePageRuleDelete,

		Schema: map[string]*schema.Schema{
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"matches": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"actions": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
			},

			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},

			"status": &schema.Schema{
				Default:  "active",
				Optional: true,
				Type:     schema.TypeString,
			},

			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCloudFlarePageRuleAction(d *schema.ResourceData) cloudflare.PageRule {
	var matchRule cloudflare.PageRuleTarget
	matchRule.Target = "url"
	matchRule.Constraint.Operator = "matches"
	matchRule.Constraint.Value = d.Get("matches").(string)

	newPageRule := cloudflare.PageRule{
		Targets: []cloudflare.PageRuleTarget{matchRule},      //d.Get("targets").(string),
		Actions: buildCloudFlarePageRuleActionsFromConfig(d), //d.Get("actions").(string),
	}

	if priority, ok := d.GetOk("priority"); ok {
		newPageRule.Priority = cloudflare.MaybeInt(priority.(int))
	}

	if status, ok := d.GetOk("status"); ok {
		newPageRule.Status = status.(string)
	}
	return newPageRule
}

func buildCloudFlarePageRuleActionsFromConfig(d *schema.ResourceData) []cloudflare.PageRuleAction {
	actions := d.Get("actions").(map[string]interface{})
	specs := make([]cloudflare.PageRuleAction, len(actions))

	i := 0
	for id, value := range actions {

		specs[i].ID = id
		specs[i].Value = value

		// Try parsing to see if it's a number since TF thinks everything is a string
		number, err := strconv.Atoi(value.(string))
		if err == nil {
			specs[i].Value = number
		}
		i++
	}

	return specs
}

func resourceCloudFlarePageRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	newPageRule := buildCloudFlarePageRuleAction(d)

	zoneName := d.Get("domain").(string)

	zoneId, err := client.ZoneIDByName(zoneName)
	if err != nil {
		return fmt.Errorf("Error finding zone %q: %s", zoneName, err)
	}

	d.Set("zone_id", zoneId)

	log.Printf("[DEBUG] CloudFlare PageRule create configuration: %#v", newPageRule)

	rule, err := client.CreatePageRule(zoneId, newPageRule)
	if err != nil {
		return fmt.Errorf("Failed to create record: %s", err)
	}

	d.SetId(rule.ID)
	// d.Set("name", record.Name)
	// d.Set("dev_mode", record.DevMode)
	// d.Set("created_on", record.CreatedOn)
	log.Printf("[DEBUG] CloudFlare PageRule ID: %#v", rule.ID)

	return resourceCloudFlarePageRuleRead(d, meta)
}

func resourceCloudFlarePageRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	rule, err := client.PageRule(d.Get("zone_id").(string), d.Id())
	if err != nil {
		return err
	}

	d.SetId(rule.ID)

	return nil
}

func resourceCloudFlarePageRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	newPageRule := buildCloudFlarePageRuleAction(d)
	rule, err := client.ChangePageRule(d.Get("zone_id").(string), d.Id(), newPageRule)
	log.Printf("[DEBUG] Updating page rule %s", d.Id())
	if err != nil {
		log.Printf("DYING")
		return err
	}

	d.SetId(rule.ID)

	return resourceCloudFlarePageRuleRead(d, meta)
}

func resourceCloudFlarePageRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)
	domain := d.Get("domain").(string)

	log.Printf("[INFO] Deleting CloudFlare PageRule: %s, %s", domain, d.Id())

	err := client.DeletePageRule(d.Get("zone_id").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting CloudFlare PageRule: %s", err)
	}

	return nil
}
