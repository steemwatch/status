package mongodb

import (
	"time"

	"github.com/steemwatch/status/checks"
	"gopkg.in/mgo.v2"
)

type pingChecker struct {
	mongoURL string
}

func newPingChecker(mongoURL string) *pingChecker {
	return &pingChecker{mongoURL}
}

func (check *pingChecker) Description() string {
	return "Connect to MongoDB"
}

func (check *pingChecker) Check(interruptCh <-chan struct{}) (*checks.CheckSummary, error) {
	session, err := mgo.DialWithTimeout(check.mongoURL, 3*time.Second)
	if err == nil {
		session.Close()
		return &checks.CheckSummary{
			Result: checks.CheckResultPassing,
		}, nil
	} else {
		return &checks.CheckSummary{
			Result:  checks.CheckResultFailing,
			Details: err.Error(),
		}, nil
	}
}

func (check *pingChecker) Period() time.Duration {
	return 5 * time.Minute
}
