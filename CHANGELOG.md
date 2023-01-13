## master (Unreleased)

NEW FEATURES:

* Complete all primitive top-level block types ([#59](https://github.com/minamijoyo/tfedit/pull/59))
* Rename references for website_domain and website_endpoint ([#60](https://github.com/minamijoyo/tfedit/pull/60))

ENHANCEMENTS:

* Update Go to v1.19 ([#55](https://github.com/minamijoyo/tfedit/pull/55))
* Update Terraform to v1.3.6 ([#56](https://github.com/minamijoyo/tfedit/pull/56))
* Update localstack to v1.3.1 ([#58](https://github.com/minamijoyo/tfedit/pull/58))

## 0.2.0 (2022/12/19)

BREAKING CHANGES:

* Redesigning the interface as a library ([#54](https://github.com/minamijoyo/tfedit/pull/54))

NEW FEATURES:

* Add support for provider meta argument ([#51](https://github.com/minamijoyo/tfedit/pull/51))
* Rename s3_force_path_style to s3_use_path_style in provider aws block ([#52](https://github.com/minamijoyo/tfedit/pull/52))
* Add DataSource type to tfwrite package ([#53](https://github.com/minamijoyo/tfedit/pull/53))

## 0.1.3 (2022/08/12)

ENHANCEMENTS:

* Use GitHub App token for updating brew formula on release ([#46](https://github.com/minamijoyo/tfedit/pull/46))

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
