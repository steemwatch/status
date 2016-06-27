package checks

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/tomb.v2"
)

var ErrInterrupted = errors.New("interrupted")

type Section struct {
	Heading  string
	Checkers []Checker
}

type Runner struct {
	sections []*Section
	results  [][]*CheckSummary

	statusCh chan<- CheckSummary

	m *sync.RWMutex
	t *tomb.Tomb
}

type RunnerOption func(*Runner)

func NewRunner(opts ...RunnerOption) *Runner {
	runner := &Runner{
		m: &sync.RWMutex{},
		t: &tomb.Tomb{},
	}

	for _, opt := range opts {
		opt(runner)
	}

	return runner
}

func SetStatusChannel(statusCh chan<- CheckSummary) RunnerOption {
	return func(runner *Runner) {
		runner.statusCh = statusCh
	}
}

func (runner *Runner) AddSection(section *Section) {
	runner.sections = append(runner.sections, section)

	summary := &CheckSummary{
		Result:    CheckResultPending,
		Details:   "Check in progress...",
		Timestamp: time.Now(),
	}

	results := make([]*CheckSummary, 0, len(section.Checkers))
	for _ = range section.Checkers {
		results = append(results, summary)
	}
	runner.results = append(runner.results, results)
}

func (runner *Runner) Start() {
	for i, section := range runner.sections {
		for j := range section.Checkers {
			runner.startChecker(i, j)
		}
	}
}

func (runner *Runner) Interrupt() {
	runner.t.Kill(nil)
}

func (runner *Runner) Wait() error {
	return runner.t.Wait()
}

func (runner *Runner) startChecker(sectionIndex, checkerIndex int) {
	runner.t.Go(func() error {
		return runner.checkLoop(sectionIndex, checkerIndex)
	})
}

func (runner *Runner) checkLoop(sectionIndex, checkerIndex int) error {
	checker := runner.sections[sectionIndex].Checkers[checkerIndex]

	for {
		runner.m.RLock()
		last := runner.results[sectionIndex][checkerIndex]
		runner.m.RUnlock()

		current, err := checker.Check(runner.t.Dying())
		if err != nil {
			if err == ErrInterrupted {
				return nil
			}
			return err
		}
		current.Description = checker.Description()
		current.Timestamp = time.Now()

		runner.m.Lock()
		runner.results[sectionIndex][checkerIndex] = current
		runner.m.Unlock()

		if runner.statusCh != nil {
			skip := last.Result == CheckResultPending && current.Result == CheckResultPassing
			skip = skip || last.Result == current.Result

			if !skip {
				select {
				case runner.statusCh <- *current:
				case <-time.After(checker.Period()):
					continue
				case <-runner.t.Dying():
					return nil
				}
			}
		}

		select {
		case <-time.After(checker.Period()):
		case <-runner.t.Dying():
			return nil
		}
	}
}

type SectionSummary struct {
	Heading string
	Checks  []*CheckSummary
}

func (runner *Runner) Results() []*SectionSummary {
	runner.m.RLock()
	defer runner.m.RUnlock()

	results := make([]*SectionSummary, 0, len(runner.sections))
	for i, section := range runner.sections {
		results = append(results, &SectionSummary{
			Heading: section.Heading,
			Checks:  runner.results[i],
		})
	}
	return results
}
