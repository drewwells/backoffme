package backoffme

type Backoffer interface {
	Retry() <-chan struct{}
	Reset()
}

type JitterType uint8

const (
	NO_JITTER JitterType = iota
	FULL_JITTER
	EQUAL_JITTER
	DECORRELATED_JITTER
)
