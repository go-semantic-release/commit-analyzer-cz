package analyzer

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

var CAVERSION = "dev"

type DefaultCommitAnalyzer struct {
	majorReleaseRules releaseRules
	minorReleaseRules releaseRules
	patchReleaseRules releaseRules
}

func (da *DefaultCommitAnalyzer) Init(m map[string]string) error {
	var err error
	da.majorReleaseRules, err = parseRules(cmp.Or(m["major_release_rules"], defaultMajorReleaseRules))
	if err != nil {
		return fmt.Errorf("failed to parse major release rules: %w", err)
	}
	da.minorReleaseRules, err = parseRules(cmp.Or(m["minor_release_rules"], defaultMinorReleaseRules))
	if err != nil {
		return fmt.Errorf("failed to parse minor release rules: %w", err)
	}
	da.patchReleaseRules, err = parseRules(cmp.Or(m["patch_release_rules"], defaultPatchReleaseRules))
	if err != nil {
		return fmt.Errorf("failed to parse patch release rules: %w", err)
	}
	return nil
}

func (da *DefaultCommitAnalyzer) Name() string {
	return "default"
}

func (da *DefaultCommitAnalyzer) Version() string {
	return CAVERSION
}

func (da *DefaultCommitAnalyzer) setTypeAndChange(c *semrel.Commit) {
	pc := parseCommit(c.Raw[0])
	if pc == nil {
		return
	}

	c.Type = pc.Type
	c.Scope = pc.Scope
	c.Message = pc.Message

	c.Change = &semrel.Change{
		// either uses the `!` convention or has a breaking change section
		Major: da.majorReleaseRules.Matches(pc) || matchesBreakingPattern(c),
		Minor: da.minorReleaseRules.Matches(pc),
		Patch: da.patchReleaseRules.Matches(pc),
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
