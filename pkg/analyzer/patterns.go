package analyzer

import (
	"regexp"
	"strings"

	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

var (
	releaseRulePattern     = regexp.MustCompile(`^([\w-\*]+)(?:\(([^\)]*)\))?(\S*)$`)
	commitPattern          = regexp.MustCompile(`^([\w-]+)(?:\(([^\)]*)\))?(\S*)\: (.*)$`)
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

func matchesBreakingPattern(c *semrel.Commit) bool {
	return breakingPattern.MatchString(strings.Join(c.Raw, "\n"))
}
