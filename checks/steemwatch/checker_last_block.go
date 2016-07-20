package steemwatch

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/steemwatch/status/checks"

	"github.com/pkg/errors"
	"github.com/tchap/steemwatch/server/routes/api/v1/info"
	"github.com/valyala/fasthttp"
)

type lastBlockChecker struct {
	infoURL string
}

func mustNewLastBlockChecker(steemwatchURL string) *lastBlockChecker {
	u, err := url.Parse(steemwatchURL)
	if err != nil {
		panic(err)
	}
	u.Path = "/api/v1/info"
	return &lastBlockChecker{u.String()}
}

func (checker *lastBlockChecker) Description() string {
	return "Check the last processed block timestamp"
}

func (checker *lastBlockChecker) Check(interruptCh <-chan struct{}) (*checks.CheckSummary, error) {
	info, err := checker.getInfo(interruptCh)
	if err != nil {
		return &checks.CheckSummary{
			Result:  checks.CheckResultFailing,
			Details: err.Error(),
		}, nil
	}

	if info.LastBlockTimestamp == nil {
		return &checks.CheckSummary{
			Result:  checks.CheckResultPending,
			Details: "last block timestamp not available yet",
		}, nil
	}

	var result checks.CheckResult
	delta := time.Now().Sub(*info.LastBlockTimestamp)
	if delta <= 10*time.Minute {
		result = checks.CheckResultPassing
	} else {
		result = checks.CheckResultFailing
	}

	details := fmt.Sprintf("next block is %v; last block processed %v ago",
		info.NextBlockNumber, delta)

	return &checks.CheckSummary{
		Result:  result,
		Details: details,
	}, nil
}

func (checker *lastBlockChecker) getInfo(interruptCh <-chan struct{}) (*info.Info, error) {
	statusCode, body, err := fasthttp.GetTimeout(nil, checker.infoURL, 10*time.Second)
	if err != nil {
		return nil, err
	} else if statusCode < 200 || statusCode >= 300 {
		return nil, errors.Errorf("HTTP response: %v %v", statusCode, string(body))
	}

	var i info.Info
	if err := json.Unmarshal(body, &i); err != nil {
		return nil, err
	}
	return &i, nil
}

func (checker *lastBlockChecker) Period() time.Duration {
	return 5 * time.Minute
}
