package backoffme

import (
	"context"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var TimeoutErr = errors.New("retry timeout")

// WaitForParams defines the parameters for the WaitFor function
type WaitForParams struct {
	// InitialDelay waits for elapsed duration before calling fn
	InitialDelay time.Duration
	// MaxDelay is the maximum duration for the exponential backoff
	MaxDelay time.Duration
	// Timeout caps the waiting to this absolute duration. Waits longer
	// than Timeout will fail with TimeoutError
	Timeout time.Duration
	// TimeoutError is the error presented when a timeout occurs
	TimeoutError error
	// JitterType adds randomness to the backoffs NO_JITTER, FULL_JITTER, EQUAL_JITTER
	JitterType JitterType
	Logger     *log.Entry
}

// NewStopRetryErr implements StopRetry interface.
func NewStopRetryErr(err error) error {
	return &stopRetryErr{msg: err}
}

// stopRetryErr is an error type that tells an exponential backoff waiter to
// stop waiting and immediately exit.
type stopRetryErr struct {
	msg error
}

// StopRetry stops a retryer from retrying.
type StopRetry interface {
	StopRetry()
}

func (s *stopRetryErr) Error() string {
	return s.msg.Error()
}

func (s *stopRetryErr) StopRetry() {
	// do nothing, implement StopRetry interface
}

// WaitFor is a generic backoff waiter that takes a context, params, and a check function.
// It uses the check function to check the condition.  If there are no errors,
// the waiter is done.  If there are errors, backoff waiting is performed depending
// on the intial delay and max delay parameters.  When the Timeout is reached,
// the TimeoutError is returned.  If the check function returns a StopRetryErr,
// then the waiter immediately returns.
func WaitFor(ctx context.Context, params WaitForParams, fn func() error) error {
	waiter := NewExpBackoff(params.InitialDelay, params.MaxDelay, params.JitterType)
	defer waiter.Reset()
	after := time.After(params.Timeout)
	var err error
	for {
		select {
		case <-waiter.Retry():
			if err = fn(); err != nil {
				if _, ok := err.(StopRetry); ok {
					return err
				}
				if params.Logger != nil {
					params.Logger.WithError(err).Warn("Waiter check did not succeed - waiting for next retry")
				}
				continue
			}
			return nil
		case <-after:
			if params.TimeoutError != nil {
				return params.TimeoutError
			}
			if err != nil {
				return errors.Wrap(err, TimeoutErr.Error())
			}
			return TimeoutErr
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
