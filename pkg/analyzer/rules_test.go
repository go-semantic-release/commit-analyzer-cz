package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRule(t *testing.T) {
	testCases := []struct {
		rule   string
		wanted *releaseRule
	}{
		{
			rule:   "feat",
			wanted: &releaseRule{Type: "feat", Scope: "*", Modifier: ""},
		},
		{
			rule:   "feat(api)",
			wanted: &releaseRule{Type: "feat", Scope: "api", Modifier: ""},
		},
		{
			rule:   "feat(*)!",
			wanted: &releaseRule{Type: "feat", Scope: "*", Modifier: "!"},
		},
		{
			rule:   "feat(api)!",
			wanted: &releaseRule{Type: "feat", Scope: "api", Modifier: "!"},
		},
		{
			rule:   "*(*)!",
			wanted: &releaseRule{Type: "*", Scope: "*", Modifier: "!"},
		},
		{
			rule:   "*(*)*",
			wanted: &releaseRule{Type: "*", Scope: "*", Modifier: "*"},
		},
		{
			rule:   "*",
			wanted: &releaseRule{Type: "*", Scope: "*", Modifier: ""},
		},
		{
			rule:   "*!",
			wanted: &releaseRule{Type: "*", Scope: "*", Modifier: "!"},
		},
		{
			rule:   "x!",
			wanted: &releaseRule{Type: "x", Scope: "*", Modifier: "!"},
		},
		{
			rule:   "x🦄",
			wanted: &releaseRule{Type: "x", Scope: "*", Modifier: "🦄"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.rule, func(t *testing.T) {
			r, err := parseRule(tc.rule)
			require.NoError(t, err)
			require.Equal(t, tc.wanted, r)
		})
	}
}
