package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {

	// The actual provider
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITLAB_TOKEN", nil),
				Description: descriptions["token"],
			},
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GITLAB_BASE_URL", ""),
				Description: descriptions["base_url"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zeus_gitlab_group":        resourceGitlabGroup(),
			"zeus_gitlab_project":      resourceGitlabProject(),
			"zeus_gitlab_project_hook": resourceGitlabProjectHook(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"token":    "The OAuth token used to connect to GitLab.",
		"base_url": "The GitLab Base API URL.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	fmt.Println("being called", d.Get("token").(string), d.Get("base_url").(string))
	config := Config{
		Token:   d.Get("token").(string),
		BaseURL: d.Get("base_url").(string),
	}

	return config.Client()
}
