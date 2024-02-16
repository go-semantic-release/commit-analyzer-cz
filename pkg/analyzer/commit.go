package analyzer

import "strings"

type parsedCommit struct {
	Type     string
	Scope    string
	Modifier string
	Message  string
}

func parseCommit(msg string) *parsedCommit {
	found := commitPattern.FindAllStringSubmatch(msg, -1)
	if len(found) < 1 {
		// commit message does not match pattern
		return nil
	}

	return &parsedCommit{
		Type:     strings.ToLower(found[0][1]),
		Scope:    found[0][2],
		Modifier: found[0][3],
		Message:  found[0][4],
	}
}
