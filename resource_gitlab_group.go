package main

import (
	"fmt"
	"log"

	"github.com/dihedron/terraform/helper/validation"
	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabGroup() *schema.Resource {
	return &schema.Resource{
		Exists: resourceGitlabGroupExists,
		Create: resourceGitlabGroupCreate,
		Read:   resourceGitlabGroupRead,
		Update: resourceGitlabGroupUpdate,
		Delete: resourceGitlabGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateName,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validatePath,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lfs_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"request_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"visibility_level": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "internal", "public"}, true),
			},
		},
	}
}

func resourceGitlabGroupSetToState(d *schema.ResourceData, group *gitlab.Group) {
	d.Set("name", group.Name)
	d.Set("path", group.Path)
	d.Set("description", group.Description)
	d.Set("lfs_enabled", group.LFSEnabled)
	d.Set("request_access_enabled", group.RequestAccessEnabled)
	d.Set("visibility_level", visibilityLevelToString(group.VisibilityLevel))
}

func resourceGitlabGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*gitlab.Client)
	group, _, err := client.Groups.GetGroup(d.Id)
	if group != nil && err == nil {
		return true, nil
	}
	return false, err
}

func resourceGitlabGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.CreateGroupOptions{
		Name: gitlab.String(d.Get("name").(string)),
		Path: gitlab.String(d.Get("path").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("lfs_enabled"); ok {
		options.LFSEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("request_access_enabled"); ok {
		options.RequestAccessEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("visibility_level"); ok {
		options.VisibilityLevel = stringToVisibilityLevel(v.(string))
	}

	log.Printf("[DEBUG] create gitlab group %q (path: %q)", options.Name, options.Path)

	group, _, err := client.Groups.CreateGroup(options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", group.ID))

	return resourceGitlabGroupRead(d, meta)
}

func resourceGitlabGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab group %s", d.Id())

	group, response, err := client.Groups.GetGroup(d.Id())
	if err != nil {
		if response.StatusCode == 404 {
			log.Printf("[WARN] removing group %s from state because it no longer exists in gitlab", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	resourceGitlabGroupSetToState(d, group)
	return nil
}

func resourceGitlabGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	options := &gitlab.UpdateGroupOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("path") {
		options.Path = gitlab.String(d.Get("path").(string))
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	if d.HasChange("lfs_enabled") {
		options.LFSEnabled = gitlab.Bool(d.Get("lfs_enabled").(bool))
	}

	if d.HasChange("request_access_enabled") {
		options.RequestAccessEnabled = gitlab.Bool(d.Get("request_access_enabled").(bool))
	}

	if d.HasChange("visibility_level") {
		options.VisibilityLevel = stringToVisibilityLevel(d.Get("visibility_level").(string))
	}

	log.Printf("[DEBUG] update gitlab group %s", d.Id())

	_, _, err := client.Groups.UpdateGroup(d.Id(), options)
	if err != nil {
		return err
	}

	return resourceGitlabGroupRead(d, meta)
}

func resourceGitlabGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] Delete gitlab group %s", d.Id())

	_, err := client.Groups.DeleteGroup(d.Id())
	return err
}
