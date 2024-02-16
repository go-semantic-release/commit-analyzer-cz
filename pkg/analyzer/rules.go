package analyzer

import (
	"cmp"
	"fmt"
	"strings"
)

var (
	defaultMajorReleaseRules = "*(*)!"
	defaultMinorReleaseRules = "feat"
	defaultPatchReleaseRules = "fix"
)

type releaseRule struct {
	Type     string
	Scope    string
	Modifier string
}

func (r *releaseRule) String() string {
	return fmt.Sprintf("%s(%s)%s", r.Type, r.Scope, r.Modifier)
}

func (r *releaseRule) Matches(commit *parsedCommit) bool {
	return (r.Type == "*" || r.Type == commit.Type) &&
		(r.Scope == "*" || r.Scope == commit.Scope) &&
		(r.Modifier == "*" || r.Modifier == commit.Modifier)
}

func parseRule(rule string) (*releaseRule, error) {
	foundRule := releaseRulePattern.FindAllStringSubmatch(rule, -1)
	if len(foundRule) < 1 {
		return nil, fmt.Errorf("cannot parse rule: %s", rule)
	}
	return &releaseRule{
		Type: strings.ToLower(foundRule[0][1]),
		// undefined scope defaults to *
		Scope:    cmp.Or(foundRule[0][2], "*"),
		Modifier: foundRule[0][3],
	}, nil
}

type releaseRules []*releaseRule

func (r releaseRules) String() string {
	ret := make([]string, len(r))
	for i, rule := range r {
		ret[i] = rule.String()
	}
	return strings.Join(ret, ",")
}

func (r releaseRules) Matches(commit *parsedCommit) bool {
	for _, rule := range r {
		if rule.Matches(commit) {
			return true
		}
	}
	return false
}

func parseRules(rules string) (releaseRules, error) {
	if rules == "" {
		return nil, fmt.Errorf("no rules provided")
	}
	ruleStrings := strings.Split(rules, ",")
	ret := make(releaseRules, len(ruleStrings))
	for i, r := range ruleStrings {
		parsed, err := parseRule(r)
		if err != nil {
			return nil, err
		}
		ret[i] = parsed
	}
	return ret, nil
}
