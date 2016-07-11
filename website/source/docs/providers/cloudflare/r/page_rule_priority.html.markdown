---
layout: "cloudflare"
page_title: "CloudFlare: cloudflare_record"
sidebar_current: "docs-cloudflare-resource-page-rule-priority"
description: |-
  Provides a Cloudflare resource to manage Page Rule priority.
---

# cloudflare\_page\_rule\_priority

Provides a Cloudflare resource to manage Page Rule priority.

## Example Usage

```
# Define the priority of CloudFlare Page Rules
resource "cloudflare_page_rule_priority" "kb_priority" {
	domain = "${var.cloudflare_domain}"
    rules = [
        "${cloudflare_page_rule.pr-everything.id}",
        "${cloudflare_page_rule.pr-other-thing.id}",
        "${cloudflare_page_rule.pr-apple.id}",
    ]
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain to add the record to
* `rules` - (Required) A list of Page Rule resource IDs. The page rules will be prioritized in a top-down manner.

## Attributes Reference

The following attributes are exported:

* `id` - The record ID
* `name` - The name of the record
* `value` - The value of the record
* `type` - The type of the record
* `ttl` - The TTL of the record
* `priority` - The priority of the record
* `hostname` - The FQDN of the record
* `proxied` - (Optional) Whether the record gets CloudFlares origin protection.

