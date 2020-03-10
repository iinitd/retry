package retry

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSimpleRetry(t *testing.T) {
	rc := New(3, 3)
	var res int
	rc.OnRun(func() error {
		i, err := mockFn(90)
		res = i
		return err
	})
	err := rc.Run()
	if err != nil {
		logf("final error %+v", err)
	}
	logf("final result %d", res)
}

func TestRetry(t *testing.T) {
	rc := New(3, 3)
	var res int
	rc.OnRun(func() error {
		i, err := mockFn(50)
		res = i
		return err
	})
	rc.OnRecord(func(currRound int, currCost float64, currErr error) {
		logf("timecost %f", currCost)
		isSucc := currErr == nil
		if currRound > 0 {
			monitor(0, 1) // 重试
			logf("retry %d %t", currRound, isSucc)
		}
	})
	err := rc.Run()
	if err != nil {
		monitor(0, 1) // 最终失败
		fmt.Println("final error", err)
		// return
	}
	fmt.Println("final result", res, rc.Costs)
}

// stub utils

func mockFn(max int) (int, error) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	i := r.Intn(100)
	if i < max {
		return i, fmt.Errorf("err %d < %d", i, max)
	}
	return i, nil
}

func logf(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}

func monitor(id int, incr int) {
	// pass
}
