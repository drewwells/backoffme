package backoffme

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestWaitFor_no_error(t *testing.T) {
	ctx := context.TODO()
	ch := make(chan error)
	eErr := errors.New("timeout error")

	go func() {
		err := WaitFor(ctx, WaitForParams{
			MaxDelay:     25 * time.Millisecond,
			Timeout:      50 * time.Millisecond,
			TimeoutError: eErr,
			JitterType:   FULL_JITTER,
		}, func() error {
			return nil
		})
		ch <- err
	}()

	if err := <-ch; err != nil {
		t.Errorf("unexpected err: %s", err)
	}

}

func TestWaitFor_timeout(t *testing.T) {
	t.Skip("waitfor does not respect max timeout if fn is still running")
	ctx := context.TODO()
	ch := make(chan error)
	eErr := errors.New("timeout error")

	go func() {
		err := WaitFor(ctx, WaitForParams{
			MaxDelay:     250 * time.Millisecond,
			Timeout:      time.Second,
			TimeoutError: eErr,
			JitterType:   FULL_JITTER,
		}, func() error {
			time.Sleep(10 * time.Second)
			return nil
		})
		ch <- err
	}()

	if err := <-ch; err != eErr {
		t.Errorf("got: %s wanted: %s", err, eErr)
	}
}

func TestWaitFor_StopEarly(t *testing.T) {
	ctx := context.TODO()
	eErr := errors.New("timeout error")
	stopErr := NewStopRetryErr(errors.New("stop"))

	var counter int

	err := WaitFor(ctx, WaitForParams{
		InitialDelay: 10 * time.Millisecond,
		MaxDelay:     2 * time.Second,
		Timeout:      10 * time.Second,
		TimeoutError: eErr,
		JitterType:   FULL_JITTER,
	}, func() error {
		counter++
		if counter > 3 {
			return stopErr
		}
		return fmt.Errorf("keep on going")
	})

	if err == nil {
		t.Fatalf("Expected stop err, got nil")
	}
	if err != stopErr {
		t.Fatalf("Expected %v, got %v", stopErr, err)
	}
	if counter != 4 {
		t.Fatalf("Expected counter = 4, got %d", counter)
	}
}
