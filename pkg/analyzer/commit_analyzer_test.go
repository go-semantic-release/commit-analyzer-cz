package analyzer

import (
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
	require.NoError(t, defaultAnalyzer.Init(map[string]string{}))
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
			&semrel.Change{Major: true, Minor: false, Patch: false},
		},
		{
			createRawCommit("f", "fix!: fixed a typo"),
			"fix",
			"",
			&semrel.Change{Major: true, Minor: false, Patch: false},
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
	require.NoError(t, defaultAnalyzer.Init(map[string]string{}))
	for _, tc := range testCases {
		t.Run(tc.RawCommit.RawMessage, func(t *testing.T) {
			analyzedCommit := defaultAnalyzer.analyzeSingleCommit(tc.RawCommit)
			require.Equal(t, tc.Type, analyzedCommit.Type, "Type")
			require.Equal(t, tc.Scope, analyzedCommit.Scope, "Scope")
			require.Equal(t, tc.Change.Major, analyzedCommit.Change.Major, "Major")
			require.Equal(t, tc.Change.Minor, analyzedCommit.Change.Minor, "Minor")
			require.Equal(t, tc.Change.Patch, analyzedCommit.Change.Patch, "Patch")
		})
	}
}

func TestReleaseRules(t *testing.T) {
	type commits []struct {
		RawCommit *semrel.RawCommit
		Change    *semrel.Change
	}
	testCases := []struct {
		Config  map[string]string
		Commits commits
	}{
		{
			Config: map[string]string{
				"major_release_rules": "feat",
				"minor_release_rules": "feat",
				"patch_release_rules": "feat",
			},
			Commits: commits{
				{
					RawCommit: createRawCommit("a", "feat: new feature"),
					Change:    &semrel.Change{Major: true, Minor: true, Patch: true},
				},
				{
					RawCommit: createRawCommit("a", "docs: new feature"),
					Change:    &semrel.Change{Major: false, Minor: false, Patch: false},
				},
				{
					RawCommit: createRawCommit("a", "feat!: new feature"),
					Change:    &semrel.Change{Major: false, Minor: false, Patch: false},
				},
				{
					RawCommit: createRawCommit("a", "feat(api): new feature"),
					Change:    &semrel.Change{Major: true, Minor: true, Patch: true},
				},
				{
					RawCommit: createRawCommit("a", "feat(api)!: new feature"),
					Change:    &semrel.Change{Major: false, Minor: false, Patch: false},
				},
			},
		},
		{
			Config: map[string]string{
				"major_release_rules": "*!",
				"minor_release_rules": "feat,chore(deps)",
				"patch_release_rules": "fix",
			},
			Commits: commits{
				{
					RawCommit: createRawCommit("a", "feat: new feature"),
					Change:    &semrel.Change{Major: false, Minor: true, Patch: false},
				},
				{
					RawCommit: createRawCommit("a", "docs: new feature"),
					Change:    &semrel.Change{Major: false, Minor: false, Patch: false},
				},
				{
					RawCommit: createRawCommit("a", "docs!: new feature"),
					Change:    &semrel.Change{Major: true, Minor: false, Patch: false},
				},
				{
					RawCommit: createRawCommit("a", "chore(deps): update dependencies"),
					Change:    &semrel.Change{Major: false, Minor: true, Patch: false},
				},
				{
					RawCommit: createRawCommit("a", "chore: cleanup"),
					Change:    &semrel.Change{Major: false, Minor: false, Patch: false},
				},
				{
					RawCommit: createRawCommit("a", "fix: bug #123"),
					Change:    &semrel.Change{Major: false, Minor: false, Patch: true},
				},
				{
					RawCommit: createRawCommit("a", "fix!: bug #123"),
					Change:    &semrel.Change{Major: true, Minor: false, Patch: false},
				},
			},
		},
	}
	for _, tc := range testCases {
		defaultAnalyzer := &DefaultCommitAnalyzer{}
		require.NoError(t, defaultAnalyzer.Init(tc.Config))
		for _, commit := range tc.Commits {
			t.Run(commit.RawCommit.RawMessage, func(t *testing.T) {
				analyzedCommit := defaultAnalyzer.analyzeSingleCommit(commit.RawCommit)
				require.Equal(t, commit.Change.Major, analyzedCommit.Change.Major, "Major")
				require.Equal(t, commit.Change.Minor, analyzedCommit.Change.Minor, "Minor")
				require.Equal(t, commit.Change.Patch, analyzedCommit.Change.Patch, "Patch")
			})
		}

	}
}
