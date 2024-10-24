package purl

import (
	"testing"
)

func TestBatchPurlLookup(t *testing.T) {
	tests := []struct {
		name     string
		purls    []string
		expected map[string]string
	}{
		{
			name: "Valid PURLs",
			purls: []string{
				"pkg:npm/%40colors/colors@1.5.0",
				"pkg:nuget/castle.core@5.1.1",
			},
			expected: map[string]string{
				"pkg:npm/%40colors/colors@1.5.0": "https://github.com/DABH/colors.js",
				"pkg:nuget/castle.core@5.1.1":    "https://github.com/castleproject/Core",
			},
		},
		{
			name: "Invalid PURL",
			purls: []string{
				"pkg:invalid/invalid@0.0.0",
			},
			expected: map[string]string{
				"pkg:invalid/invalid@0.0.0": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := BatchPurlLookup(tt.purls)
			if err != nil {
				t.Fatalf("BatchPurlLookup returned an error: %v", err)
			}

			for _, purl := range tt.purls {
				expectedURL := tt.expected[purl]
				if url, exists := results[purl]; !exists || url != expectedURL {
					t.Errorf("For purl %s, expected URL %s, but got %s", purl, expectedURL, url)
				} else {
					t.Logf("Purl: %s, GitHub URL: %s", purl, url)
				}
			}
		})
	}
}
