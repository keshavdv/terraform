---
layout: "cloudflare"
page_title: "CloudFlare: cloudflare_zone"
sidebar_current: "docs-cloudflare-resource-zone"
description: |-
  Provides a Cloudflare zone resource.
---

# cloudflare\_zone

Provides a Cloudflare zone resource.

## Example Usage

```
# Add a domain to CloudFlare
resource "cloudflare_zone" "foobar" {
	domain = "mywebsite.com"
 	organization = "My Org"
   	plan = "Enterprise Website"
    settings {
		always_online = "off"
		browser_cache_ttl = 14400
		browser_check = "off"
		cache_level = "basic"
		challenge_ttl = 2600
		email_obfuscation = "off"
		max_upload = 500
		server_side_exclude = "off"
		minify {
			html = "on"
			js = "off"
			css = "on"
		}
   	}
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain to add the zone to
* `plan` - (Required) The name of the plan
* `organization` - (Required) The name of the organization to which this domain should be added
* `jump_start` - (Optional) Automatically attempt to fetch existing DNS zones (default: true)
* `settings` - (Optional) A map of additional zone settings defined by the [offical API](https://api.cloudflare.com/#zone-settings-properties).

    **Note:** Setting the `security_header` option is currently unsupported

## Attributes Reference

The following attributes are exported:

* `id` - The zone ID
* `name` - The domain name
* `dev_mode` - The interval (in seconds) from when development mode expires (positive integer) or last expired (negative integer) for the domain. If development mode has never been enabled, this value is 0.
* `created_on` - When the zone was created
