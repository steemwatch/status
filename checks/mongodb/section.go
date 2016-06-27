package mongodb

import (
	"github.com/steemwatch/status/checks"
)

func NewSection(mongoURL string) *checks.Section {
	return &checks.Section{
		Heading: "MongoDB",
		Checkers: []checks.Checker{
			newPingChecker(mongoURL),
		},
	}
}
