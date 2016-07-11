package cloudflare

import (
	"fmt"
	"log"
	"reflect"

	"github.com/cloudflare/cloudflare-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCloudFlareZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFlareZoneCreate,
		Read:   resourceCloudFlareZoneRead,
		Update: resourceCloudFlareZoneUpdate,
		Delete: resourceCloudFlareZoneDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"organization": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"plan": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"settings": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advanced_ddos": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"always_online": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"browser_cache_ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"browser_check": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"cache_level": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"challenge_ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"cname_flattening": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"development_mode": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"edge_cache_ttl": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"email_obfuscation": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"hotlink_protection": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"http2": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_geolocation": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ipv6": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"max_upload": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"minify": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"css": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"html": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"js": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"mirage": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"mobile_redirect": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"status": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"mobile_subdomain": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"strip_uri": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"origin_error_page_pass_thru": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"polish": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"prefetch_preload": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"pseudo_ipv4": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"response_buffering": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"rocket_loader": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"security_level": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"server_side_exclude": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sha1_support": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"sort_query_string_for_cache": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"ssl": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"tls_1_2_only": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"tls_client_auth": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"true_client_ip_header": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"waf": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"websockets": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"dev_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_on": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"jump_start": &schema.Schema{
				Default:  true,
				Optional: true,
				Type:     schema.TypeBool,
			},
		},
	}
}

func resourceCloudFlareZoneCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	var matchingOrg cloudflare.Organization
	if d.Get("organization") != nil {
		userDetails, err := client.UserDetails()
		if err != nil {
			return fmt.Errorf("Error getting user's organizations: %s", err)
		}
		found := false
		expected_org := d.Get("organization").(string)
		for _, org := range userDetails.Organizations {
			if org.Name == expected_org {
				log.Printf("[DEBUG] Found organization matching '%s'", expected_org)
				found = true
				matchingOrg = org
				break
			}
		}

		if !found {
			return fmt.Errorf("Could not find organization '%s'", expected_org)
		}
	}

	zone, err := client.CreateZone(d.Get("name").(string), d.Get("jump_start").(bool), matchingOrg)
	if err != nil {
		return fmt.Errorf("Failed to create zone: %s", err)
	}

	log.Printf("[DEBUG] CloudFlare Zone created for '%s' with ID '%s", zone.Name, zone.ID)

	if d.Get("plan") != nil {
		plans, err := client.AvailableZonePlans(zone.ID)
		if err != nil {
			return fmt.Errorf("Error getting available plans: %s", err)
		}

		found := false
		expected_plan := d.Get("plan").(string)
		for _, plan := range plans {
			if plan.Name == expected_plan {
				log.Printf("[DEBUG] Found plan matching '%s'", expected_plan)
				found = true

				_, err := client.ZoneSetPlan(zone.ID, plan)
				if err != nil {
					return fmt.Errorf("Failed to set plan for zone '%s': %s", d.Get("domain"), err)
				}
			}
		}

		if !found {
			return fmt.Errorf("Could not find plan '%s'", expected_plan)
		}
	}

	d.SetId(zone.ID)

	return nil
	// resourceCloudFlareZoneUpdate(d, meta)
}

func zoneSettingsToMap(settings []cloudflare.ZoneSetting) map[string]interface{} {
	settings_map := make(map[string]interface{}, len(settings))

	for _, setting := range settings {
		log.Printf("[DEBUG] %s is %s", setting.ID, setting.Value)

		if setting.ID == "mobile_redirect" || setting.ID == "security_header" || setting.ID == "minify" {
			// settings_map[setting.ID] = []interface{}{setting.Value}
		} else {
			settings_map[setting.ID] = setting.Value
		}
	}

	return settings_map
}

func buildCloudFlareZoneSettings(d *schema.ResourceData, meta interface{}) ([]cloudflare.ZoneSetting, error) {
	var settings_to_apply []cloudflare.ZoneSetting

	if v, ok := d.GetOk("settings"); ok {
		settings := v.([]interface{})
		if len(settings) > 1 {
			return nil, fmt.Errorf("You can only define settings once per zone")
		}

		for k, v := range settings[0].(map[string]interface{}) {
			if _, ok := d.GetOk("settings.0." + k); ok {
				setting := cloudflare.ZoneSetting{}
				setting.ID = k
				setting.Value = v
				log.Printf("[DEBUG] CloudFlare setting '%s' is '%v'", k, v)
				log.Printf("[DEBUG] %s", reflect.TypeOf(v))
				settings_to_apply = append(settings_to_apply, setting)
			}
		}
	}
	log.Printf("[DEBUG] CloudFlare settings ares '%v'", settings_to_apply)

	return settings_to_apply, nil
}

func resourceCloudFlareZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	zone, err := client.ZoneDetails("/" + d.Id())
	if err != nil {
		return err
	}

	d.SetId(zone.ID)
	d.Set("name", zone.Name)
	d.Set("dev_mode", zone.DevMode)
	d.Set("created_on", zone.CreatedOn)

	// if v, ok := d.GetOk("settings"); ok {
	// 	expectedSettings := v.(*schema.Set).List()[0]
	// 	settings, err := client.GetZoneSettings(d.Id())
	// 	if err != nil {
	// 		return err
	// 	}
	// 	log.Printf("map is %s", zoneSettingsToMap(d, settings))

	// 	foo := zoneSettingsToMap(settings)
	// 	// err = d.Set("settings", []interface{}{foo})

	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// resourceCloudFlareZoneUpdate(d, meta)
	return nil
}

func resourceCloudFlareZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	// return nil
	client := meta.(*cloudflare.API)
	d.Partial(true)

	settings_to_apply, err := buildCloudFlareZoneSettings(d, meta)
	if err != nil {
		return fmt.Errorf("Failed to build zone settings: %s", err)
	}

	_, err = client.EditZoneSettings(d.Id(), settings_to_apply)
	if err != nil {
		return fmt.Errorf("Failed to edit zone settings: %s", err)
	}
	d.Partial(false)

	return resourceCloudFlareZoneRead(d, meta)
}

func resourceCloudFlareZoneDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*cloudflare.API)

	_, err := client.DeleteZone("/" + d.Id())
	if err != nil {
		return fmt.Errorf("Failed to delete zone: %s", err)
	}

	log.Printf("[DEBUG] CloudFlare Zone '%s' deleted", d.Get("name"))

	d.SetId("")

	return nil
}
