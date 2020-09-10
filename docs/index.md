# Artifactory Provider

The [Artifactory](https://jfrog.com/artifactory/) provider is used to interact with the
resources supported by Artifactory. The provider needs to be configured
with the proper credentials before it can be used.

- Available Resources
    * [Groups](./resources/artifactory_group.md)
    * [Users](./resources/artifactory_user.md)
    * [Permission Target](./resources/artifactory_permission_target.md)
    * [Local Repositories](./resources/artifactory_local_repository.md)
    * [Remote Repositories](./resources/artifactory_remote_repository.md)
    * [Replication Configurations](./resources/artifactory_replication_config.md)
    * [Single Replication Configurations](./resources/artifactory_single_replication_config.md)
    * [Virtual Repositories](./resources/artifactory_virtual_repository.md)
    * [Certificates](./resources/artifactory_certificate.md)

- Available Datasources
    * [File](./data-sources/artifactory_file.md)
    * [FileInfo](./data-sources/artifactory_fileinfo.md)

- Deprecated Resources
    * [Permission Targets (V1 API)](./resources/artifactory_permission_target_v1.md)

## Example Usage
```hcl
# Required for [Terraform 0.13](https://www.terraform.io/upgrade-guides/0-13.html)
terraform {
  required_providers {
    artifactory = {
      source  = "terraform.example.com/atlassian/artifactory"
      version = "2.0.0"
    }
  }
}

# Configure the Artifactory provider
provider "artifactory" {
  url = "${var.artifactory_url}"
  username = "${var.artifactory_username}"
  password = "${var.artifactory_password}"
}

# Create a new repository
resource "artifactory_local_repository" "pypi-libs" {
  key             = "pypi-libs"
  package_type    = "pypi"
  repo_layout_ref = "simple-default"
  description     = "A pypi repository for python packages"
}
```

## Authentication
The Artifactory provider supports multiple means of authentication. The following methods are supported:
    * Basic Auth
    * Bearer Token
    * JFrog API Key Header

### Basic Auth
Basic auth may be used by adding a `username` and `password` field to the provider block
Getting this value from the environment is supported with the `ARTIFACTORY_USERNAME` and `ARITFACTORY_PASSWORD` variable

Usage:
```hcl
# Configure the Artifactory provider
provider "artifactory" {
  url = "artifactory.site.com"
  username = "myusername"
  password = "mypassword"
}
```

### Bearer Token
Artifactory access tokens may be used via the Authorization header by providing the `access_token` field to the provider
block. Getting this value from the environment is supported with the `ARTIFACTORY_ACCESS_TOKEN` variable

Usage:
```hcl
# Configure the Artifactory provider
provider "artifactory" {
  url = "artifactory.site.com"
  access_token = "abc...xy"
}
```

### JFrog API Key Header
Artifactory API keys may be used via the `X-JFrog-Art-Api` header by providing the `api_key` field in the provider block.
Getting this value from the environment is supported with the `ARTIFACTORY_API_KEY` variable

Usage:
```hcl
# Configure the Artifactory provider
provider "artifactory" {
  url = "artifactory.site.com"
  api_key = "abc...xy"
}
```

## Argument Reference

The following arguments are supported:

* `url` - (Required) URL of Artifactory. This can also be sourced from the `ARTIFACTORY_URL` environment variable.
* `username` - (Optional) Username for basic auth. Requires `password` to be set. 
    Conflicts with `api_key`, and `access_token`. This can also be sourced from the `ARTIFACTORY_USERNAME` environment variable.
* `password` - (Optional) Password for basic auth. Requires `username` to be set. 
    Conflicts with `api_key`, and `access_token`. This can also be sourced from the `ARTIFACTORY_PASSWORD` environment variable.
* `api_key` - (Optional) API key for api auth. Uses `X-JFrog-Art-Api` header. 
    Conflicts with `username`, `password`, and `access_token`. This can also be sourced from the `ARTIFACTORY_API_KEY` environment variable.
* `access_token` - (Optional) API key for token auth. Uses `Authorization: Bearer` header. 
    Conflicts with `username` and `password`, and `api_key`. This can also be sourced from the `ARTIFACTORY_ACCESS_TOKEN` environment variable.