package backoffme

import (
	"math"
	"math/rand"
	"sync"
	"time"
)

type expBackoff struct {
	d, max time.Duration
	tries  uint32
	lock   sync.Mutex
	j      JitterType
}

func NewExpBackoff(delay, max time.Duration, jitter JitterType) Backoffer {
	if delay == 0 {
		delay = 20 * time.Millisecond
	}
	return &expBackoff{
		d:    delay,
		max:  max,
		lock: sync.Mutex{},
		j:    jitter,
	}
}

func (e *expBackoff) Retry() <-chan struct{} {
	e.lock.Lock()
	wait := GetDelay(e.tries, e.d, e.max, e.j)
	c := make(chan struct{})
	go func() {
		defer e.lock.Unlock()
		defer func() { e.tries++ }()
		select {
		case <-time.After(wait):
			close(c)
		}
	}()
	return c
}

func GetDelay(tries uint32, initDelay, maxDelay time.Duration, jitter JitterType) time.Duration {
	fMaxDelay := float64(maxDelay)
	wait := float64(initDelay) * math.Pow(2, float64(tries))
	if wait > fMaxDelay {
		return maxDelay
	}
	switch jitter {
	case FULL_JITTER:
		wait = float64(time.Duration(rand.Int63n(int64(wait))))
	case EQUAL_JITTER:
		iWait := int64(wait)
		// integer division, but who cares
		wait = float64(time.Duration(rand.Int63n(iWait/2) + iWait/2))
	case NO_JITTER:
	}
	if wait > fMaxDelay {
		wait = fMaxDelay
	}

	return time.Duration(wait)
}

func (e *expBackoff) Reset() {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.tries = 0
}
