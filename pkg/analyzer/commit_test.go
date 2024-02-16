package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommit(t *testing.T) {
	testCases := []struct {
		message string
		wanted  *parsedCommit
	}{
		{
			message: "feat: new feature",
			wanted:  &parsedCommit{"feat", "", "", "new feature"},
		},
		{
			message: "feat!: new feature",
			wanted:  &parsedCommit{"feat", "", "!", "new feature"},
		},
		{
			message: "feat(api): new feature",
			wanted:  &parsedCommit{"feat", "api", "", "new feature"},
		},
		{
			message: "feat(api): a(b): c:",
			wanted:  &parsedCommit{"feat", "api", "", "a(b): c:"},
		},
		{
			message: "feat(new cool-api): feature",
			wanted:  &parsedCommit{"feat", "new cool-api", "", "feature"},
		},
		{
			message: "feat(😅): cool",
			wanted:  &parsedCommit{"feat", "😅", "", "cool"},
		},
		{
			message: "this-is-also(valid): cool",
			wanted:  &parsedCommit{"this-is-also", "valid", "", "cool"},
		},
		{
			message: "feat((x)): test",
			wanted:  &parsedCommit{"feat", "(x", ")", "test"},
		},
		{
			message: "feat(x)?!: test",
			wanted:  &parsedCommit{"feat", "x", "?!", "test"},
		},
		{
			message: "feat(x): test",
			wanted:  &parsedCommit{"feat", "x", "", "test"},
		},
		{
			message: "feat(x): : test",
			wanted:  &parsedCommit{"feat", "x", "", ": test"},
		},
		{
			message: "feat!: test",
			wanted:  &parsedCommit{"feat", "", "!", "test"},
		},
		// invalid messages
		{
			message: "feat (new api): feature",
			wanted:  nil,
		},
		{
			message: "feat:test",
			wanted:  nil,
		},
		{
			message: "🚀(🦄): emojis!",
			wanted:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.message, func(t *testing.T) {
			c := parseCommit(tc.message)
			require.Equal(t, tc.wanted, c)
		})
	}
}
