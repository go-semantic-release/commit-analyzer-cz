package analyzer

import (
	"strings"

	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

var CAVERSION = "dev"

type DefaultCommitAnalyzer struct{}

func (da *DefaultCommitAnalyzer) Init(m map[string]string) error {
	// TODO: implement config parsing
	return nil
}

func (da *DefaultCommitAnalyzer) Name() string {
	return "default"
}

func (da *DefaultCommitAnalyzer) Version() string {
	return CAVERSION
}

func (da *DefaultCommitAnalyzer) setTypeAndChange(c *semrel.Commit) {
	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		// commit message does not match pattern
		return
	}

	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	c.Message = found[0][4]

	c.Change = &semrel.Change{
		// either uses the `!` convention or has a breaking change section
		Major: found[0][3] == "!" || matchesBreakingPattern(c),
		Minor: c.Type == "feat",
		Patch: c.Type == "fix",
	}
}

func (da *DefaultCommitAnalyzer) analyzeSingleCommit(rawCommit *semrel.RawCommit) *semrel.Commit {
	c := &semrel.Commit{
		SHA:         rawCommit.SHA,
		Raw:         strings.Split(rawCommit.RawMessage, "\n"),
		Change:      &semrel.Change{},
		Annotations: rawCommit.Annotations,
	}
	c.Annotations["mentioned_issues"] = extractMentions(mentionedIssuesPattern, rawCommit.RawMessage)
	c.Annotations["mentioned_users"] = extractMentions(mentionedUsersPattern, rawCommit.RawMessage)

	da.setTypeAndChange(c)
	return c
}

func (da *DefaultCommitAnalyzer) Analyze(rawCommits []*semrel.RawCommit) []*semrel.Commit {
	ret := make([]*semrel.Commit, len(rawCommits))
	for i, c := range rawCommits {
		ret[i] = da.analyzeSingleCommit(c)
	}
	return ret
}
