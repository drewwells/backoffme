package backoffme

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func Test_ExpBackoff_NO_JITTER(t *testing.T) {
	t.Skip()
	b := NewExpBackoff(100*time.Millisecond, 10*time.Second, NO_JITTER)
	waits := []time.Duration{}
	for i := 0; i < 5; i++ {
		start := time.Now()
		<-b.Retry()
		waits = append(waits, time.Since(start))
		t.Log("waited", waits[i])
	}
	for i := 0; i < len(waits)-1; i++ {
		delta := waits[i]*2 - waits[i+1]
		if delta > 500*time.Millisecond {
			t.Fatal("Delta was too big:", delta)
		}
	}

	b = NewExpBackoff(100*time.Millisecond, 1*time.Second, NO_JITTER)
	waits = []time.Duration{}
	for i := 0; i < 10; i++ {
		start := time.Now()
		<-b.Retry()
		waits = append(waits, time.Since(start))
		//fmt.Println("waited", waits[i])
	}
	for i := 0; i < len(waits)-7; i++ {
		delta := waits[i]*2 - waits[i+1]
		if delta > 500*time.Millisecond {
			t.Fatal("Delta was too big:", delta, "waits:", waits)
		}
	}
	for i := 4; i < len(waits); i++ {
		delta := waits[i] - 1*time.Second
		if delta > 500*time.Millisecond {
			t.Fatal("Wait should have been close to 1s, was too far away:", delta)
		}
	}
}

func Test_ExpBackoff_FULL_JITTER(t *testing.T) {
	t.Skip()
	b := NewExpBackoff(100*time.Millisecond, 10*time.Second, FULL_JITTER)
	fmt.Println(runtime.NumGoroutine(), "goroutines")
	for j := 0; j < 2; j++ {
		for i := 0; i < 5; i++ {
			start := time.Now()
			select {
			case <-b.Retry():
			}
			//t.Log("Waited", time.Since(start))
			fmt.Println("waited", time.Since(start))
		}
		b.Reset()
	}
}

func Test_ExpBackoff_EQUAL_JITTER(t *testing.T) {
	t.Skip()
	b := NewExpBackoff(100*time.Millisecond, 10*time.Second, EQUAL_JITTER)
	fmt.Println(runtime.NumGoroutine(), "goroutines")
	for j := 0; j < 2; j++ {
		for i := 0; i < 5; i++ {
			start := time.Now()
			select {
			case <-b.Retry():
			}
			//t.Log("Waited", time.Since(start))
			fmt.Println("waited", time.Since(start))
		}
		b.Reset()
	}
}
