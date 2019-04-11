package io

import (
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"testing"
	"time"
)

const (
	pollInterval = 500 * time.Millisecond
)

// Exists returns true if the given file exists, false otherwise.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// WaitExists blocks and waits until the given file exists.
func WaitExists(t *testing.T, name string, timeout time.Duration) error {
	t.Helper()
	t.Logf("Waiting for file %s to become available", name)
	defer t.Logf("Done waiting for file %s to become available", name)

	first := true
	return wait.Poll(pollInterval, timeout, func() (bool, error) {
		if !first {
			t.Logf("Waiting for file %s to become available", name)
		}
		first = false

		exists := Exists(name)
		return exists, nil

	})
}

