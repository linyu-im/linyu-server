package timer

import "time"

type TaskFunc func(interface{}) bool

func Timer(delay time.Duration, fun TaskFunc, param interface{}) {
	go func() {
		if fun == nil {
			return
		}
		t := time.NewTimer(delay)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				if fun(param) == false {
					return
				}
				t.Reset(delay)
			}
		}
	}()
}
