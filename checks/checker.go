package checks

import (
	"time"
)

type CheckResult string

const (
	CheckResultPending CheckResult = "Pending"
	CheckResultPassing CheckResult = "Passing"
	CheckResultFailing CheckResult = "Failing"
)

type CheckSummary struct {
	Description string
	Result      CheckResult
	Details     string
	Timestamp   time.Time
}

type Checker interface {
	Description() string
	Check(interruptCh <-chan struct{}) (*CheckSummary, error)
	Period() time.Duration
}
