package steemwatch

import (
	"time"

	"github.com/steemwatch/status/checks"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type connectionChecker struct {
	url string
}

func newConnectionChecker(steemwatchURL string) *connectionChecker {
	return &connectionChecker{steemwatchURL}
}

func (checker *connectionChecker) Description() string {
	return "GET steemwatch.com home page (no proxy involved)"
}

func (checker *connectionChecker) Check(interruptCh <-chan struct{}) (*checks.CheckSummary, error) {
	err := checker.doCheck(interruptCh)
	if err == nil {
		return &checks.CheckSummary{
			Result:  checks.CheckResultPassing,
			Details: "HTTP response: 200 OK",
		}, nil
	} else {
		return &checks.CheckSummary{
			Result:  checks.CheckResultFailing,
			Details: err.Error(),
		}, nil
	}
}

func (checker *connectionChecker) doCheck(interruptCh <-chan struct{}) error {
	statusCode, body, err := fasthttp.GetTimeout(nil, checker.url, 10*time.Second)
	if err == nil {
		if statusCode < 200 || statusCode >= 300 {
			err = errors.Errorf("HTTP response: %v %v", statusCode, string(body))
		}

	}
	return err
}

func (checker *connectionChecker) Period() time.Duration {
	return 5 * time.Minute
}
