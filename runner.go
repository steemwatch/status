package main

import (
	"github.com/steemwatch/status/checks"
	"github.com/steemwatch/status/checks/mongodb"
	"github.com/steemwatch/status/checks/steem"
	"github.com/steemwatch/status/checks/steemwatch"
)

func startRunner(opts ...checks.RunnerOption) *checks.Runner {
	runner := checks.NewRunner(opts...)
	runner.AddSection(steemwatch.NewSection("http://localhost:8080"))
	runner.AddSection(steem.NewSection("ws://localhost:8090"))
	runner.AddSection(mongodb.NewSection("mongodb://localhost/steemwatch"))
	runner.Start()
	return runner
}
