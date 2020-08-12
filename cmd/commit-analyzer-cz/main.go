package main

import (
	defaultAnalyzer "github.com/go-semantic-release/commit-analyzer-cz/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/analyzer"
	"github.com/go-semantic-release/semantic-release/v2/pkg/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		CommitAnalyzer: func() analyzer.CommitAnalyzer {
			return &defaultAnalyzer.DefaultCommitAnalyzer{}
		},
	})
}
