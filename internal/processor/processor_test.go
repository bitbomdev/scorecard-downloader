package processor

import (
	"testing"
)

func TestParseGitHubPURL(t *testing.T) {
	tests := []struct {
		input         string
		expectedOwner string
		expectedRepo  string
	}{
		{"github.com/owner/repo", "owner", "repo"},
		{"github.com/owner/repo/extra", "owner", "repo"},
		{"github.com/owner", "", ""},
		{"notgithub.com/owner/repo", "", ""},
		{"github.com/", "", ""},
		{"", "", ""},
		{"https://github.com/DABH/colors.js", "DABH", "colors.js"},
		{"https://github.com/castleproject/Core", "castleproject", "Core"},
	}

	for _, test := range tests {
		owner, repo := parseGitHubPURL(test.input)
		if owner != test.expectedOwner || repo != test.expectedRepo {
			t.Errorf("parseGitHubPURL(%q) = (%q, %q); want (%q, %q)", test.input, owner, repo, test.expectedOwner, test.expectedRepo)
		}
	}
}
