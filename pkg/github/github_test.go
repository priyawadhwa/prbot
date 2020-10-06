package github

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestPRContainsLabel(t *testing.T) {
	tests := []struct {
		description string
		labels      []*github.Label
		required    []string
		expected    bool
	}{
		{
			description: "all labels exists",
			labels: []*github.Label{
				githubLabel("one"),
				githubLabel("two"),
				githubLabel("three"),
			},
			required: []string{"one", "two"},
			expected: true,
		},
		{
			description: "all labels don't exist",
			labels: []*github.Label{
				githubLabel("one"),
				githubLabel("two"),
				githubLabel("three"),
			},
			required: []string{"one", "two", "four"},
		},
		{
			description: "no labels required",
			labels: []*github.Label{
				githubLabel("one"),
				githubLabel("two"),
				githubLabel("three"),
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			actual := prContainsLabels(test.labels, test.required)
			if actual != test.expected {
				t.Fatalf("expected %v got %v", test.expected, actual)
			}
		})
	}
}

func githubLabel(name string) *github.Label {
	return &github.Label{
		Name: &name,
	}
}
