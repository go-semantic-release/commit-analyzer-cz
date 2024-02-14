package analyzer

import (
	"regexp"
	"strings"

	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

var (
	CAVERSION              = "dev"
	commitPattern          = regexp.MustCompile(`^([^\s\(\!]+)(?:\(([^\)]*)\))?(\!)?\: (.*)$`)
	breakingPattern        = regexp.MustCompile("BREAKING CHANGES?")
	mentionedIssuesPattern = regexp.MustCompile(`#(\d+)`)
	mentionedUsersPattern  = regexp.MustCompile(`(?i)@([a-z\d]([a-z\d]|-[a-z\d])+)`)
)

func extractMentions(re *regexp.Regexp, s string) string {
	ret := make([]string, 0)
	for _, m := range re.FindAllStringSubmatch(s, -1) {
		ret = append(ret, m[1])
	}
	return strings.Join(ret, ",")
}

type DefaultCommitAnalyzer struct{}

func (da *DefaultCommitAnalyzer) Init(_ map[string]string) error {
	return nil
}

func (da *DefaultCommitAnalyzer) Name() string {
	return "default"
}

func (da *DefaultCommitAnalyzer) Version() string {
	return CAVERSION
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

	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	breakingChange := found[0][3]
	c.Message = found[0][4]

	isMajorChange := breakingPattern.MatchString(rawCommit.RawMessage)
	isMinorChange := c.Type == "feat"
	isPatchChange := c.Type == "fix"

	if len(breakingChange) > 0 {
		isMajorChange = true
		isMinorChange = false
		isPatchChange = false
	}

	c.Change = &semrel.Change{
		Major: isMajorChange,
		Minor: isMinorChange,
		Patch: isPatchChange,
	}
	return c
}

func (da *DefaultCommitAnalyzer) Analyze(rawCommits []*semrel.RawCommit) []*semrel.Commit {
	ret := make([]*semrel.Commit, len(rawCommits))
	for i, c := range rawCommits {
		ret[i] = da.analyzeSingleCommit(c)
	}
	return ret
}
