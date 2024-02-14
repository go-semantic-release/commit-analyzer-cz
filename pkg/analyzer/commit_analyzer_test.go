package analyzer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
	"github.com/stretchr/testify/require"
)

func createRawCommit(sha, message string) *semrel.RawCommit {
	return &semrel.RawCommit{
		SHA:        sha,
		RawMessage: message,
		Annotations: map[string]string{
			"author_name": "test",
		},
	}
}

func TestAnnotations(t *testing.T) {
	defaultAnalyzer := &DefaultCommitAnalyzer{}
	rawCommit := createRawCommit("a", "fix: bug #123 and #243\nthanks @Test-user for providing this fix\n\nCloses #22")
	commit := defaultAnalyzer.analyzeSingleCommit(rawCommit)
	require.Equal(t, rawCommit.SHA, commit.SHA)
	require.Equal(t, rawCommit.RawMessage, strings.Join(commit.Raw, "\n"))
	require.Equal(t, "test", commit.Annotations["author_name"])
	require.Equal(t, "123,243,22", commit.Annotations["mentioned_issues"])
	require.Equal(t, "Test-user", commit.Annotations["mentioned_users"])
}

func TestDefaultAnalyzer(t *testing.T) {
	testCases := []struct {
		RawCommit *semrel.RawCommit
		Type      string
		Scope     string
		Change    *semrel.Change
	}{
		{
			createRawCommit("a", "feat: new feature"),
			"feat",
			"",
			&semrel.Change{Major: false, Minor: true, Patch: false},
		},
		{
			createRawCommit("b", "feat(web): new feature"),
			"feat",
			"web",
			&semrel.Change{Major: false, Minor: true, Patch: false},
		},
		{
			createRawCommit("c", "new feature"),
			"",
			"",
			&semrel.Change{Major: false, Minor: false, Patch: false},
		},
		{
			createRawCommit("d", "chore: break\nBREAKING CHANGE: breaks everything"),
			"chore",
			"",
			&semrel.Change{Major: true, Minor: false, Patch: false},
		},
		{
			createRawCommit("e", "feat!: modified login endpoint"),
			"feat",
			"",
			&semrel.Change{Major: true, Minor: true, Patch: false},
		},
		{
			createRawCommit("f", "fix!: fixed a typo"),
			"fix",
			"",
			&semrel.Change{Major: true, Minor: false, Patch: true},
		},
		{
			createRawCommit("g", "refactor(parser)!: drop support for Node 6\n\nBREAKING CHANGE: refactor to use JavaScript features not available in Node 6."),
			"refactor",
			"parser",
			&semrel.Change{Major: true, Minor: false, Patch: false},
		},
		{
			createRawCommit("h", "docs: added more documentation"),
			"docs",
			"",
			&semrel.Change{Major: false, Minor: false, Patch: false},
		},
		{
			createRawCommit("i", "chore: moved README.md to root"),
			"chore",
			"",
			&semrel.Change{Major: false, Minor: false, Patch: false},
		},
		{
			createRawCommit("i", "feat(deps): update deps\n\nBREAKING CHANGE: update to new version of dep"),
			"feat",
			"deps",
			&semrel.Change{Major: true, Minor: true, Patch: false},
		},
	}

	defaultAnalyzer := &DefaultCommitAnalyzer{}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("AnalyzeCommitMessage: %s", tc.RawCommit.RawMessage), func(t *testing.T) {
			analyzedCommit := defaultAnalyzer.analyzeSingleCommit(tc.RawCommit)
			require.Equal(t, tc.Type, analyzedCommit.Type, "Type")
			require.Equal(t, tc.Scope, analyzedCommit.Scope, "Scope")
			require.Equal(t, tc.Change.Major, analyzedCommit.Change.Major, "Major")
			require.Equal(t, tc.Change.Minor, analyzedCommit.Change.Minor, "Minor")
			require.Equal(t, tc.Change.Patch, analyzedCommit.Change.Patch, "Patch")
		})
	}
}

func TestCommitPattern(t *testing.T) {
	testCases := []struct {
		message string
		wanted  []string
	}{
		{
			message: "feat: new feature",
			wanted:  []string{"feat", "", "", "new feature"},
		},
		{
			message: "feat!: new feature",
			wanted:  []string{"feat", "", "!", "new feature"},
		},
		{
			message: "feat(api): new feature",
			wanted:  []string{"feat", "api", "", "new feature"},
		},
		{
			message: "feat(api): a(b): c:",
			wanted:  []string{"feat", "api", "", "a(b): c:"},
		},
		{
			message: "feat(new cool-api): feature",
			wanted:  []string{"feat", "new cool-api", "", "feature"},
		},
		{
			message: "feat(😅): cool",
			wanted:  []string{"feat", "😅", "", "cool"},
		},
		{
			message: "this-is-also(valid): cool",
			wanted:  []string{"this-is-also", "valid", "", "cool"},
		},
		{
			message: "🚀(🦄): emojis!",
			wanted:  []string{"🚀", "🦄", "", "emojis!"},
		},
		// invalid messages
		{
			message: "feat (new api): feature",
			wanted:  nil,
		},
		{
			message: "feat((x)): test",
			wanted:  nil,
		},
		{
			message: "feat:test",
			wanted:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("CommitPattern: %s", tc.message), func(t *testing.T) {
			found := commitPattern.FindAllStringSubmatch(tc.message, -1)
			if len(tc.wanted) == 0 {
				require.Len(t, found, 0)
				return
			}
			require.Len(t, found, 1)
			require.Len(t, found[0], 5)
			require.Equal(t, tc.wanted, found[0][1:])
		})
	}
}
