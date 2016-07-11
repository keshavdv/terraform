---
layout: "cloudflare"
page_title: "CloudFlare: cloudflare_pagerule"
sidebar_current: "docs-cloudflare-resource-page-rule"
description: |-
  Provides a Cloudflare Page Rule resource.
---

# cloudflare\_page\_rule

Provides a Cloudflare Page Rule resource.

## Example Usage

```
# Add a Page Rule to the domain
resource "cloudflare_page_rule" "pr-cache" {
	domain = "${var.cloudflare_domain}"
    matches = "${var.cloudflare_domain}/images/*"
    actions {
        browser_cache_ttl = 1800
        cache_level = "cache_everything"
        edge_cache_ttl = 600
    }
    status = "active"
}
```

## Argument Reference

Most of these arguments directly correspond to the
[offical API](https://api.cloudflare.com/#page-rules-for-a-zone-properties).


* `domain` - (Required) The domain to add the Page Rule to
* `matches` - (Required) The url to match
* `actions` - (Required) The set of actions to perform for matching requests. See the CloudFlare API documentations for possible options.
* `status` - (Optional) The status of the page rule (one of: active, paused)
* `priority` - (Optional) The priority of the Page Rule

    **Note:** If you are managing multiple Page Rules via Terraform, the priority field may not work as expected unless you explicitly define a dependency between resources. Alternatively, use the [`cloudflare_pagerule_priority`](page_rule_priority.html) resource instead.

## Attributes Reference

The following attributes are exported:

* `id` - The Page Rule ID
* `zone_id` - Zone identifier tag
* `priority` - The priority of the Page Rule
* `status` - The status of the Page Rule

