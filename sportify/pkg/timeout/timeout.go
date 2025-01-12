package timeout

import (
	"fmt"
	"time"
)

var ErrTimeout = fmt.Errorf("Timeout")

// Timeout executes the function f with a timeout.
func Timeout(timeout time.Duration, f func() error) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- f()
	}()
	select {
	case <-time.After(timeout):
		return fmt.Errorf("%w: after %v", ErrTimeout, timeout)
	case err := <-errChan:
		return err
	}
}
