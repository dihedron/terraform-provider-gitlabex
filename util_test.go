package main

import (
	"testing"

	"github.com/xanzy/go-gitlab"
)

func TestGitlab_validation(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "invalid",
			ErrCount: 1,
		},
		{
			Value:    "valid_one",
			ErrCount: 0,
		},
		{
			Value:    "valid_two",
			ErrCount: 0,
		},
	}

	validationFunc := validateValueFunc([]string{"valid_one", "valid_two"})

	for _, tc := range cases {
		_, errors := validationFunc(tc.Value, "test_arg")

		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected 1 validation error")
		}
	}
}

func TestGitlab_validateName(t *testing.T) {
	// A Group/Project name can contain only letters, digits, '_', '.', dash and
	// space.
	cases := []struct {
		String string
		Errors int
	}{
		{
			String: "My New App 02",
			Errors: 0,
		},
		{
			String: "My New App 02#",
			Errors: 1,
		},
		{
			String: "My New App - 02.",
			Errors: 0,
		},
		{
			String: "My New App - 02;",
			Errors: 1,
		},
	}
	for _, tc := range cases {
		_, errors := validateName(tc.String, "name")
		if len(errors) != tc.Errors {
			t.Fatalf("got %d errors expected %d", len(errors), tc.Errors)
		}
	}
}

func TestGitlab_validatePath(t *testing.T) {
	// A Group/Project name can contain only letters, digits, '_', '.', dash and
	// space.
	cases := []struct {
		String string
		Errors int
	}{
		{
			String: "My New App 02",
			Errors: 1,
		},
		{
			String: "My New App 02#",
			Errors: 1,
		},
		{
			String: "My-New-App-02#",
			Errors: 1,
		},
		{
			String: "My New App - 02.",
			Errors: 1,
		},
		{
			String: "My New-App-02.",
			Errors: 1,
		},
		{
			String: "My New-App-02.atom",
			Errors: 2,
		},
		{
			String: "My New-App-02.git",
			Errors: 2,
		},
	}
	for _, tc := range cases {
		_, errors := validatePath(tc.String, "path")
		if len(errors) != tc.Errors {
			t.Fatalf("%s - got %d errors expected %d", tc.String, len(errors), tc.Errors)
		}
	}
}

func TestGitlab_visibilityHelpers(t *testing.T) {
	cases := []struct {
		String string
		Level  gitlab.VisibilityLevelValue
	}{
		{
			String: "private",
			Level:  gitlab.PrivateVisibility,
		},
		{
			String: "public",
			Level:  gitlab.PublicVisibility,
		},
	}

	for _, tc := range cases {
		level := stringToVisibilityLevel(tc.String)
		if level == nil || *level != tc.Level {
			t.Fatalf("got %v expected %v", level, tc.Level)
		}

		sv := visibilityLevelToString(tc.Level)
		if sv == nil || *sv != tc.String {
			t.Fatalf("got %v expected %v", sv, tc.String)
		}
	}
}
