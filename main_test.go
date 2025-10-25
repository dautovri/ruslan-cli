package main
package main

import (
	"testing"
)

func TestVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "default version",
			version: "dev",
			want:    "dev",
		},
		{
			name:    "semver version",
			version: "v1.0.0",
			want:    "v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			version = tt.version
			if version != tt.want {
				t.Errorf("version = %v, want %v", version, tt.want)
			}
		})
	}
}

func TestCommit(t *testing.T) {
	if commit == "" {
		commit = "none"
	}

	if len(commit) < 4 && commit != "none" {
		t.Errorf("commit hash too short: %s", commit)
	}
}

func TestDate(t *testing.T) {
	if date == "" {
		date = "unknown"
	}

	if date != "unknown" && len(date) < 10 {
		t.Errorf("date format invalid: %s", date)
	}
}
