package steem

import (
	"time"

	"github.com/go-steem/rpc"
	"github.com/steemwatch/status/checks"
)

type connectionChecker struct {
	endpointURL string
}

func newConnectionChecker(mongoURL string) *connectionChecker {
	return &connectionChecker{mongoURL}
}

func (checker *connectionChecker) Description() string {
	return "Access the private Steem RPC endpoint"
}

func (checker *connectionChecker) Check(interruptCh <-chan struct{}) (*checks.CheckSummary, error) {
	version, err := checker.doCheck(interruptCh)
	if err == nil {
		return &checks.CheckSummary{
			Result:  checks.CheckResultPassing,
			Details: "blockchain version is " + version,
		}, nil
	} else {
		return &checks.CheckSummary{
			Result:  checks.CheckResultFailing,
			Details: err.Error(),
		}, nil
	}
}

func (checker *connectionChecker) doCheck(interruptCh <-chan struct{}) (string, error) {
	client, err := rpc.Dial(checker.endpointURL)
	if err != nil {
		return "", err
	}
	defer client.Close()

	config, err := client.GetConfig()
	if err != nil {
		return "", err
	}

	return config.SteemitBlockchainVersion, nil
}

func (checker *connectionChecker) Period() time.Duration {
	return 5 * time.Minute
}
