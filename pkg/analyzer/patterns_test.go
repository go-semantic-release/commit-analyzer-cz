package analyzer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractIssues(t *testing.T) {
	testCases := []struct {
		message string
		wanted  string
	}{
		{
			message: "feat: new feature #123",
			wanted:  "123",
		},
		{
			message: "feat!: new feature closes #123 and #456",
			wanted:  "123,456",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.message, func(t *testing.T) {
			issues := extractMentions(mentionedIssuesPattern, testCase.message)
			require.Equal(t, testCase.wanted, issues)
		})
	}
}

func TestExtractMentions(t *testing.T) {
	testCases := []struct {
		message string
		wanted  string
	}{
		{
			message: "feat: new feature by @user",
			wanted:  "user",
		},
		{
			message: "feat!: new feature by @user and @user-2",
			wanted:  "user,user-2",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.message, func(t *testing.T) {
			issues := extractMentions(mentionedUsersPattern, testCase.message)
			require.Equal(t, testCase.wanted, issues)
		})
	}
}
