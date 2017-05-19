package main

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProject() *schema.Resource {
	return &schema.Resource{
		//SchemaVersion: 1,
		//MigrateState:  resourceGitlabProjectMigrateState,
		Exists: resourceGitlabProjectExists,
		Create: resourceGitlabProjectCreate,
		Read:   resourceGitlabProjectRead,
		Update: resourceGitlabProjectUpdate,
		Delete: resourceGitlabProjectDelete,

		Schema: map[string]*schema.Schema{
			// all these fieds can be set at creation/update time
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the new project.",
				Required:     true,
				ValidateFunc: validateName,
			},
			"path": {
				Type:         schema.TypeString,
				Description:  "Custom repository name for new project; by default generated based on name.",
				Optional:     true,
				Computed:     true,
				ValidateFunc: validatePath,
			},
			"default_branch": {
				Type:        schema.TypeString,
				Description: "'master' by default.",
				Optional:    true,
			},
			// the following field can be used to create the project within an
			// existing group, or to move it into it; the specified id must
			// correspond to an existing group namespace; if the namespace
			// corresponds to a user namespace, the project is moved into that
			// user's namespace; this effectively corresponsads to  call to
			// the "CreateProjectForUser()" API
			"namespace_id": {
				Type:        schema.TypeInt,
				Description: "Namespace for the new project; by default it's the current user's namespace, but it can\nbe the ID of a group.",
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Short project description.",
				Optional:    true,
			},
			"issues_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable issues for this project.",
				Optional:    true,
				Computed:    true,
			},
			"merge_requests_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable merge requests for this project.",
				Optional:    true,
				Computed:    true,
			},
			"builds_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable builds for this project.",
				Optional:    true,
				Computed:    true,
			},
			"wiki_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable builds for this project.",
				Optional:    true,
				Computed:    true,
			},
			"snippets_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable builds for this project.",
				Optional:    true,
				Computed:    true,
			},
			"container_registry_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable container registry for this project.",
				Optional:    true,
				Computed:    true,
			},
			"shared_runners_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable shared runners for this project.",
				Optional:    true,
				Computed:    true,
			},
			"visibility_level": {
				Type: schema.TypeString,
				Description: `The visbility level of the project; can be one of: 
* private  - project access must be granted explicitly for each user 
* internal - the project can be cloned by any logged in user
* public   - the project can be cloned without any authentication
`,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"private", "internal", "public"}, true),
			},
			"import_url": {
				Type:        schema.TypeString,
				Description: "URL to import repository from.",
				Optional:    true,
			},
			"public_builds": {
				Type:        schema.TypeBool,
				Description: "If true, builds can be viewed by non-project-members.",
				Optional:    true,
				Computed:    true,
			},
			"only_allow_merge_if_build_succeeds": {
				Type:        schema.TypeBool,
				Description: "Set whether merge requests can only be merged with successful builds.",
				Optional:    true,
			},
			"only_allow_merge_if_all_discussions_are_resolved": {
				Type:        schema.TypeBool,
				Description: "Set whether merge requests can only be merged when all the discussions are resolved.",
				Optional:    true,
			},
			"lfs_enabled": {
				Type:        schema.TypeBool,
				Description: "Enable Large File Support (LFS).",
				Optional:    true,
				Computed:    true,
			},
			"request_access_enabled": {
				Type:        schema.TypeBool,
				Description: "Allow users to request member access.",
				Optional:    true,
			},
			// all the following fields are computed, and are not stored in the
			// Terraform state
			"ssh_url_to_repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_url_to_repo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name_with_namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path_with_namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"open_issues_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"approvals_before_merge": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString, // formatted according to RFC3339
				Computed: true,
			},
			"last_activity_at": {
				Type:     schema.TypeString, // formatted according to RFC3339
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"archived": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"avatar_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"forks_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"stars_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"runners_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"forked_from_project_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

/*
func resourceGitlabProjectMigrateState(from int, previous *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error)) {

}
*/

func resourceGitlabProjectSetToState(d *schema.ResourceData, project *gitlab.Project) {
	d.Set("name", project.Name)
	d.Set("path", project.Path)
	d.Set("default_branch", project.DefaultBranch)
	d.Set("namespace_id", project.Namespace.ID)
	d.Set("description", project.Description)
	d.Set("issues_enabled", project.IssuesEnabled)
	d.Set("merge_requests_enabled", project.MergeRequestsEnabled)
	d.Set("builds_enabled", project.BuildsEnabled)
	d.Set("wiki_enabled", project.WikiEnabled)
	d.Set("snippets_enabled", project.SnippetsEnabled)
	d.Set("container_registry_enabled", project.ContainerRegistryEnabled)
	d.Set("shared_runners_enabled", project.SharedRunnersEnabled)
	d.Set("visibility_level", visibilityLevelToString(project.VisibilityLevel))
	// NOTE: import_url is only used at creation time
	d.Set("public_builds", project.PublicBuilds)
	d.Set("only_allow_merge_if_build_succeeds", project.OnlyAllowMergeIfBuildSucceeds)
	d.Set("only_allow_merge_if_all_discussions_are_resolved", project.OnlyAllowMergeIfAllDiscussionsAreResolved)
	d.Set("lfs_enabled", project.LFSEnabled)
	d.Set("request_access_enabled", project.RequestAccessEnabled)
}

func resourceGitlabProjectExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*gitlab.Client)
	project, _, err := client.Projects.GetProject(d.Id)
	if project != nil && err == nil {
		return true, nil
	}
	return false, err
}

func resourceGitlabProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.CreateProjectOptions{
		Name: gitlab.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("path"); ok {
		options.Path = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("default_branch"); ok {
		options.DefaultBranch = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("namespace_id"); ok {
		if _, err := checkNamespace(client, v.(int)); err != nil {
			return err
		}
		options.NamespaceID = gitlab.Int(v.(int))
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	if v, ok := d.GetOk("issues_enabled"); ok {
		options.IssuesEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("merge_requests_enabled"); ok {
		options.MergeRequestsEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("builds_enabled"); ok {
		options.BuildsEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("wiki_enabled"); ok {
		options.WikiEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("snippets_enabled"); ok {
		options.SnippetsEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("container_registry_enabled"); ok {
		options.ContainerRegistryEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("shared_runners_enabled"); ok {
		options.SharedRunnersEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("visibility_level"); ok {
		options.VisibilityLevel = stringToVisibilityLevel(v.(string))
	}

	if v, ok := d.GetOk("public_builds"); ok {
		options.PublicBuilds = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("only_allow_merge_if_build_succeeds"); ok {
		options.OnlyAllowMergeIfBuildSucceeds = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("only_allow_merge_if_all_discussions_are_resolved"); ok {
		options.OnlyAllowMergeIfAllDiscussionsAreResolved = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("lfs_enabled"); ok {
		options.LFSEnabled = gitlab.Bool(v.(bool))
	}

	if v, ok := d.GetOk("request_access_enabled"); ok {
		options.RequestAccessEnabled = gitlab.Bool(v.(bool))
	}

	log.Printf("[DEBUG] create gitlab project %q", options.Name)

	project, _, err := client.Projects.CreateProject(options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", project.ID))

	return resourceGitlabProjectRead(d, meta)
}

func resourceGitlabProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab project %s", d.Id())

	project, response, err := client.Projects.GetProject(d.Id())
	if err != nil {
		if response.StatusCode == 404 {
			log.Printf("[WARN] removing project %s from state because it no longer exists in gitlab", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	resourceGitlabProjectSetToState(d, project)
	return nil
}

func resourceGitlabProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	options := &gitlab.EditProjectOptions{}

	if d.HasChange("name") && len(d.Get("name").(string)) > 0 {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("path") && len(d.Get("path").(string)) > 0 {
		options.Path = gitlab.String(d.Get("path").(string))
	}

	// the Group's namespace id is the ID of the group itself.

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	if d.HasChange("issues_enabled") {
		options.IssuesEnabled = gitlab.Bool(d.Get("issues_enabled").(bool))
	}

	if d.HasChange("merge_requests_enabled") {
		options.MergeRequestsEnabled = gitlab.Bool(d.Get("merge_requests_enabled").(bool))
	}

	if d.HasChange("builds_enabled") {
		options.BuildsEnabled = gitlab.Bool(d.Get("builds_enabled").(bool))
	}

	if d.HasChange("wiki_enabled") {
		options.WikiEnabled = gitlab.Bool(d.Get("wiki_enabled").(bool))
	}

	if d.HasChange("snippets_enabled") {
		options.SnippetsEnabled = gitlab.Bool(d.Get("snippets_enabled").(bool))
	}

	if d.HasChange("container_registry_enabled") {
		options.ContainerRegistryEnabled = gitlab.Bool(d.Get("container_registry_enabled").(bool))
	}

	if d.HasChange("shared_runners_enabled") {
		options.SharedRunnersEnabled = gitlab.Bool(d.Get("shared_runners_enabled").(bool))
	}

	if d.HasChange("visibility_level") {
		options.VisibilityLevel = stringToVisibilityLevel(d.Get("visibility_level").(string))
	}

	if d.HasChange("public_builds") {
		options.PublicBuilds = gitlab.Bool(d.Get("public_builds").(bool))
	}

	if d.HasChange("only_allow_merge_if_build_succeeds") {
		options.OnlyAllowMergeIfBuildSucceeds = gitlab.Bool(d.Get("only_allow_merge_if_build_succeeds").(bool))
	}

	if d.HasChange("only_allow_merge_if_all_discussions_are_resolved") {
		options.OnlyAllowMergeIfAllDiscussionsAreResolved = gitlab.Bool(d.Get("only_allow_merge_if_all_discussions_are_resolved").(bool))
	}

	if d.HasChange("lfs_enabled") {
		options.LFSEnabled = gitlab.Bool(d.Get("lfs_enabled").(bool))
	}

	if d.HasChange("request_access_enabled") {
		options.RequestAccessEnabled = gitlab.Bool(d.Get("request_access_enabled").(bool))
	}

	log.Printf("[DEBUG] update gitlab project %s", d.Id())

	_, _, err := client.Projects.EditProject(d.Id(), options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectRead(d, meta)
}

func resourceGitlabProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] Delete gitlab project %s", d.Id())

	_, err := client.Projects.DeleteProject(d.Id())
	return err
}
