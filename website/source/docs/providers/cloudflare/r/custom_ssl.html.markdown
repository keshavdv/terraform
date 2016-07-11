---
layout: "cloudflare"
page_title: "CloudFlare: cloudflare_custom_ssl"
sidebar_current: "docs-cloudflare-resource-custom-ssl"
description: |-
  Provides a Cloudflare custom SSL resource
---

# cloudflare\_custom\_ssl

Provides a Cloudflare custom SSL resource.

## Example Usage

```
# Add a record to the domain
resource "cloudflare_custom_ssl" "website_ssl" {
	domain = "${var.cloudflare_domain}"
	certificate = "<PEM encoded certificate>"
	private_key = "<PEM encoded certificate>"
	bundle_method = "ubiquitous"
}
```

## Argument Reference

Most of these arguments directly correspond to the
[offical API](https://api.cloudflare.com/#custom-ssl-for-a-zone-properties).

* `domain` - (Required) The domain to add the custom SSL certificate to
* `certificate` - (Required) The certificate in PEM format
* `private_key` - (Required) The private key in PEM format
* `bundle_method` - (Optional) The desired bundle method  (valid values: ubiquitous, optimal, force)

## Attributes Reference

The following attributes are exported:

* `id` - The record ID
* `zone_id` - Zone identifier tag
* `bundle_method` - The value of the record
* `hosts` - The hosts included in the certificate
* `issuer` - The certificate authority that issued the certificate
* `signature` - The type of hash used for the certificate
* `status` - Status of the zone's custom SSL (one of: active, expired, deleted)
* `uploaded_on` - When the certificate was uploaded to CloudFlare
* `modified_on` - When the certificate was last modified
* `expires_on` - When the certificate from the authority expires
* `priority` - The order/priority in which the certificate will be used in a request. Higher numbers will be tried first.