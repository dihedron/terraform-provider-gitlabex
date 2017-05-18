package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var (
	validName = regexp.MustCompile(`^[a-zA-Z0-9_\.\- ]+$`)
	validPath = regexp.MustCompile(`^[a-zA-Z0-9_\.][a-zA-Z0-9_\.\-]*[a-zA-Z0-9_\-]+$`)
)

// A Group/Project name can contain only letters, digits, '_', '.', dash and
// space.
func validateName(v interface{}, k string) (we []string, errors []error) {
	value := v.(string)
	if !validName.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q is an invalid name: it can contain only letters, digits, '_', '.', dash and space", value))
	}
	return
}

// A Group/Project path can contain only letters, digits, '_', '-' and '.'; it
// cannot start with '-' or end in '.', '.git' or '.atom'.
func validatePath(v interface{}, k string) (we []string, errors []error) {
	value := v.(string)
	if !validPath.MatchString(value) {
		errors = append(errors, fmt.Errorf("%q is an invalid path: it can contain only letters, digits, '_', '-' and '.'; it cannot start with '-' or end in '.'", value))
	}
	if strings.HasSuffix(value, ".atom") {
		errors = append(errors, fmt.Errorf("%q is an invalid path: it cannot end in .atom", value))
	}
	if strings.HasSuffix(value, ".git") {
		errors = append(errors, fmt.Errorf("%q is an invalid path: it cannot end in .git", value))
	}
	return
}

func validateRegexpFunc(regexp string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		if !validName.MatchString(value) {
			errors = append(errors, fmt.Errorf("%s is an invalid %s", value, k))
		}
		return
	}
}

// copied from ../github/util.go
func validateValueFunc(values []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		valid := false
		for _, role := range values {
			if value == role {
				valid = true
				break
			}
		}

		if !valid {
			errors = append(errors, fmt.Errorf("%s is an invalid value for argument %s", value, k))
		}
		return
	}
}

func stringToVisibilityLevel(s string) *gitlab.VisibilityLevelValue {
	lookup := map[string]gitlab.VisibilityLevelValue{
		"private":  gitlab.PrivateVisibility,
		"internal": gitlab.InternalVisibility,
		"public":   gitlab.PublicVisibility,
	}

	value, ok := lookup[s]
	if !ok {
		return nil
	}
	return &value
}

func visibilityLevelToString(v gitlab.VisibilityLevelValue) *string {
	lookup := map[gitlab.VisibilityLevelValue]string{
		gitlab.PrivateVisibility:  "private",
		gitlab.InternalVisibility: "internal",
		gitlab.PublicVisibility:   "public",
	}
	value, ok := lookup[v]
	if !ok {
		return nil
	}
	return &value
}

// namespaces handling is a bit complex: if the ID provided corresponds
// to the id of a group namespace, then the project should be moved into
// that group; if the namespace corresponds to a user namespace, then
// the project should be moved (i.e. assigned) to that user (this
// should happen through the CreateProjectForUser API and its options); if
// the namespace ID does not xist, there is an error in the plan
func checkNamespace(client *gitlab.Client, id int) (string, error) {
	namespaces, _, err := client.Namespaces.ListNamespaces(&gitlab.ListNamespacesOptions{})

	if err != nil {
		return "", fmt.Errorf("Error getting list of namespaces: %s", err)
	}

	for _, namespace := range namespaces {
		if namespace.ID == id {
			log.Printf("[DEBUG] Namespace with ID found, of type: %s", namespace.Kind)
			return namespace.Kind, nil
		}
	}
	return "", fmt.Errorf("Invalid namespace ID: %d", id)
}
