package steem

import (
	"github.com/steemwatch/status/checks"
)

func NewSection(steemdRPCEndpointURL string) *checks.Section {
	return &checks.Section{
		Heading: "Steem Daemon",
		Checkers: []checks.Checker{
			newConnectionChecker(steemdRPCEndpointURL),
		},
	}
}
