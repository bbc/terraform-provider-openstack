module github.com/terraform-provider-openstack/terraform-provider-openstack

go 1.14

require (
	github.com/gophercloud/gophercloud v0.15.0
	github.com/gophercloud/utils v0.0.0-20210113034859-6f548432055a
	github.com/hashicorp/terraform-plugin-sdk v1.16.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/stretchr/testify v1.4.0
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/gophercloud/gophercloud => github.com/bbc/gophercloud v0.15.1-0.20210217110349-1fe61c1403ec
