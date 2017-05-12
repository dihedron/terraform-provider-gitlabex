package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGitlabGroup_basic(t *testing.T) {
	var group gitlab.Group
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			// Create a group with all the features on
			{
				Config: testAccGitlabGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("zeus_gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:        fmt.Sprintf("foo-%d", rInt),
						Path:        fmt.Sprintf("bar-%d", rInt),
						Description: "Terraform acceptance tests",
						/*
							VisibilityLevel:      20,
							LFSEnabled:           true,
							RequestAccessEnabled: true,
						*/
					}),
				),
			},
			// Update the group to turn the features off
			{
				Config: testAccGitlabGroupUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("zeus_gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:        fmt.Sprintf("foo-%d", rInt),
						Path:        fmt.Sprintf("bar-%d", rInt),
						Description: "Terraform acceptance tests!",
						/*
							VisibilityLevel: 20,
						*/
					}),
				),
			},
			//Update the group to turn the features on again
			{
				Config: testAccGitlabGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("zeus_gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:        fmt.Sprintf("foo-%d", rInt),
						Path:        fmt.Sprintf("bar-%d", rInt),
						Description: "Terraform acceptance tests",
						/*
							VisibilityLevel:      20,
							LFSEnabled:           true,
							RequestAccessEnabled: true,
						*/
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabGroupExists(n string, group *gitlab.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		repoName := rs.Primary.ID
		if repoName == "" {
			return fmt.Errorf("No group ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		gotGroup, _, err := conn.Groups.GetGroup(repoName)
		if err != nil {
			return err
		}
		*group = *gotGroup
		return nil
	}
}

type testAccGitlabGroupExpectedAttributes struct {
	Name        string
	Path        string
	Description string
	/*
		VisibilityLevel      gitlab.VisibilityLevelValue
		LFSEnabled           bool
		RequestAccessEnabled bool
	*/
}

func testAccCheckGitlabGroupAttributes(group *gitlab.Group, want *testAccGitlabGroupExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if group.Name != want.Name {
			return fmt.Errorf("got group %q; want %q", group.Name, want.Name)
		}
		if group.Name != want.Path {
			return fmt.Errorf("got path %q; want %q", group.Path, want.Path)
		}
		if group.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", group.Description, want.Description)
		}
		/*
			if group.VisibilityLevel != want.VisibilityLevel {
				return fmt.Errorf("got default branch %q; want %q", group.VisibilityLevel, want.VisibilityLevel)
			}

			if group.LFSEnabled != want.LFSEnabled {
				return fmt.Errorf("got lfs_enabled %t; want %t", group.LFSEnabled, want.LFSEnabled)
			}

			if group.RequestAccessEnabled != want.RequestAccessEnabled {
				return fmt.Errorf("got request_access_enabled %t; want %t", group.RequestAccessEnabled, want.RequestAccessEnabled)
			}
		*/
		return nil
	}
}

func testAccCheckGitlabGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zeus_gitlab_group" {
			continue
		}

		gotRepo, resp, err := conn.Groups.GetGroup(rs.Primary.ID)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				return fmt.Errorf("Repository still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabGroupConfig(rInt int) string {
	return fmt.Sprintf(`
resource "zeus_gitlab_group" "foo" {
  name = "foo-%d"
  path = "bar-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  # visibility_level = "public"
}
	`, rInt, rInt)
}

func testAccGitlabGroupUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "zeus_gitlab_group" "foo" {
  name = "foo-%d"
  path = "bar-%d"
  description = "Terraform acceptance tests!"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  # visibility_level = "public"

  #	lfs_enabled = false
  #	request_accesss_enabled = false
}
	`, rInt, rInt)
}
