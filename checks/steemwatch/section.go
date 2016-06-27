package steemwatch

import (
	"github.com/steemwatch/status/checks"
)

func NewSection(steemwatchURL string) *checks.Section {
	return &checks.Section{
		Heading: "SteemWatch",
		Checkers: []checks.Checker{
			newConnectionChecker(steemwatchURL),
		},
	}
}
