package retry

import (
	"fmt"
	"time"
)

type Runable func() error
type Recordable func(currRound int, currCost float64, currErr error)
type Errorable func()
type Costs []float64

type RetryContainer struct {
	// logger  logbale
	Rounds   int
	Costs    Costs
	Interval time.Duration
	fnRun    Runable
	fnRecord Recordable
}

func New(rounds int, interval int) *RetryContainer {
	return &RetryContainer{
		Rounds:   rounds,
		Costs:    []float64{},
		Interval: time.Duration(interval) * time.Millisecond,
	}
}

//OnRun ..
func (rc *RetryContainer) OnRun(fn Runable) {
	rc.fnRun = fn
}

//OnRetry ..
func (rc *RetryContainer) OnRecord(fn Recordable) {
	rc.fnRecord = fn
}

//Run ..
func (rc *RetryContainer) Run() error {

	if rc.fnRun == nil {
		return fmt.Errorf("retry run not init")
	}

	var (
		i   int
		err error
	)

	for ; i < rc.Rounds; i++ {
		startTime := time.Now()
		err = rc.fnRun()
		cost := time.Since(startTime).Seconds()
		if rc.fnRecord != nil {
			rc.fnRecord(i, cost, err)
		}
		rc.Costs = append(rc.Costs, cost)
		if err == nil {
			break
		}
		time.Sleep(rc.Interval)
	}

	return err
}
