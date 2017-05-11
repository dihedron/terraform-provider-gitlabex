variable "tenant" {
  type        = "string"
  default     = "tenant-x"
}

provider "gitlabex" {
#	token       		= ""
	base_url 			  = "http://localhost/api/v3/"
}

resource "gitlabex_group" "group01" {
	name = "group01"
	path ="path01"
	description="my first test group"
}

resource "gitlabex_project" "infrastructure-repo" {
	name        		= "${var.tenant}-infrastructure"
	description 		= "The infrastructure configuration GitLab repository"
	visibility_level 	= "public"
}
