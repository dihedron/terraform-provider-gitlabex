# Terraform GitLab 

A Terraform provider for the latest GitLab features, as an external plugin.
This plugin is a testbed for newer features such as Group support, nested 
Projects, code checkin; it will be adapted and pull-requested for merger
into the upstream Terraform builtin provider as soon as it is considered
sufficiently mature and stable.

## Installation

You can easily install the latest version with the following :

```
go get -u github.com/dihedron/terraform-provider-gitlabx
```

Then add the plugin to your local `.terraformrc` :

```
cat >> ~/.terraformrc <<EOF
providers {
  gitlabx = "${GOPATH}/bin/terraform-provider-gitlabx"
}
EOF
```

