package cloudflare

import (
	"crypto/sha256"
	"fmt"
	"math/rand"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDFLARE_EMAIL", nil),
				Description: "A registered CloudFlare email address.",
			},

			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUDFLARE_TOKEN", nil),
				Description: "The token key for API operations.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cloudflare_record":             resourceCloudFlareRecord(),
			"cloudflare_zone":               resourceCloudFlareZone(),
			"cloudflare_page_rule":          resourceCloudFlarePageRule(),
			"cloudflare_page_rule_priority": resourceCloudFlarePageRulePriority(),
			"cloudflare_custom_ssl":         resourceCloudFlareCustomSSL(),
			// "cloudflare_custom_ssl_priority":   resourceCloudFlareZone(),
			// "cloudflare_custom_page":   resourceCloudFlareZone(),
			// "cloudflare_firewall_rule":   resourceCloudFlareZone(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Email: d.Get("email").(string),
		Token: d.Get("token").(string),
	}

	return config.Client()
}

func randomHash() string {
	data := make([]byte, 10)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}
	return fmt.Sprintf("%x", sha256.Sum256(data))
}
