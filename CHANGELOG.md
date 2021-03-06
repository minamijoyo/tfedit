## master (Unreleased)

## 0.1.2 (2022/07/06)

BUG FIXES:

* Map mfa_delete true/false => Enabled/Disabled for aws_s3_bucket_versioning ([#41](https://github.com/minamijoyo/tfedit/pull/41))
* Suppress creating a migration file when no action ([#42](https://github.com/minamijoyo/tfedit/pull/42))
* Suppress adding invalid days_after_initiation for aws_s3_bucket_lifecycle_configuration ([#43](https://github.com/minamijoyo/tfedit/pull/43))
* Fix invalid filter and tags for aws_s3_bucket_lifecycle_configuration ([#45](https://github.com/minamijoyo/tfedit/pull/45))

ENHANCEMENTS:

* Use a native cache feature in actions/setup-go ([#44](https://github.com/minamijoyo/tfedit/pull/44))

## 0.1.1 (2022/06/16)

ENHANCEMENTS:

* Update Go to v1.17.10 and Alpine to v3.16 ([#36](https://github.com/minamijoyo/tfedit/pull/36))
* Update hcl to v2.12.0 and hcledit to v0.2.4 ([#37](https://github.com/minamijoyo/tfedit/pull/37))
* Update Go to v1.18.3 ([#38](https://github.com/minamijoyo/tfedit/pull/38))
* Update terraform-json to v0.14.0 ([#39](https://github.com/minamijoyo/tfedit/pull/39))

## 0.1.0 (2022/06/07)

ENHANCEMENTS:

* Add a note about an limitation of aws_s3_bucket.lifecycle_rule.id ([#28](https://github.com/minamijoyo/tfedit/pull/28))
* Add support for Terraform v1.2 ([#31](https://github.com/minamijoyo/tfedit/pull/31))
* Read Go version from .go-version on GitHub Actions ([#32](https://github.com/minamijoyo/tfedit/pull/32))

## 0.0.3 (2022/05/10)

NEW FEATURES:

* Add support for count and for_each ([#26](https://github.com/minamijoyo/tfedit/pull/26))

## 0.0.2 (2022/05/02)

NEW FEATURES:

* Generate a migration file for import from Terraform plan ([#25](https://github.com/minamijoyo/tfedit/pull/25))

## 0.0.1 (2022/04/15)

Initial release
