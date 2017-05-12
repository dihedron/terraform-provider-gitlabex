variable "tenant" {
  type        = "string"
  default     = "tenant-x"
}

provider "zeus_gitlab" {
#	token       		= ""
	base_url 			  = "http://localhost/api/v3/"
}

resource "zeus_gitlab_group" "group01" {
	name = "group01"
	path ="path01"
	description="my first test group"
}

resource "zeus_gitlab_project" "infrastructure-repo" {
	name        		= "${var.tenant}-infrastructure"
	path				= ${zeus_gitlab_group.group1.path}/${var.tenant}-infrastructure"
	description 		= "The infrastructure configuration GitLab repository"
	visibility_level 	= "public"
}
