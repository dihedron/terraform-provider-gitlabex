/*
 * The name of the project as a variable.
 */
variable "project" {
  type        		= "string"
  default     		= "mynewapp"
}

/*
 * The GitLab provider.
 */
provider "gitlabx" {
#	token			= "<token goes here>"
	base_url 		= "http://localhost/api/v3/"
}

/*
 * The tenant's project group.
 */
resource "gitlabx_group" "group" {
	name 				= "${var.project}"
	path 				= "${var.project}-group"
	description			= "The GitLab projects group for the ${var.project} application"
	visibility_level 	= "private"
}

/*
 * The infrastructure definitions repository.
 */ 
resource "gitlabx_project" "infrastructure" {
	name        		= "${var.project}-infrastructure"
	namespace_id		= "${gitlabx_group.group.id}"
	description 		= "The infrastructure configuration repository for the ${var.project} "
	visibility_level 	= "private"
}

/*
 * The first microservice.
 */
resource "gitlabx_project" "microservice-01" {
	name        		= "${var.project}-microservice-01"
	namespace_id		= "${gitlabx_group.group.id}"
	description 		= "The source repository for the first microservice in ${var.project}"
	visibility_level 	= "private"
}

/*
 * The second microservice.
 */
resource "gitlabx_project" "microservice-02" {
	name        		= "${var.project}-microservice-02"
	namespace_id		= "${gitlabx_group.group.id}"
	description 		= "The source repository for the second microservice in ${var.project}"
	visibility_level 	= "private"
}

/*
 * The third microservice.
 */
resource "gitlabx_project" "microservice-03" {
	name        		= "${var.project}-microservice-03"
	namespace_id		= "${gitlabx_group.group.id}"
	description 		= "The source repository for the third microservice in ${var.project}"
	visibility_level 	= "private"
}