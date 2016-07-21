package steem

import (
	"time"

	"github.com/steemwatch/status/checks"

	"github.com/go-steem/rpc"
	"github.com/go-steem/rpc/transports/websocket"
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
	t, err := websocket.NewTransport(checker.endpointURL)
	if err != nil {
		return "", err
	}

	client := rpc.NewClient(t)
	defer client.Close()

	config, err := client.Database.GetConfig()
	if err != nil {
		return "", err
	}

	return config.SteemitBlockchainVersion, nil
}

func (checker *connectionChecker) Period() time.Duration {
	return 5 * time.Minute
}
