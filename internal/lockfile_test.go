package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLockFile_TryLock(t *testing.T) {
	// Use temp directory for testing
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "test.pid")

	lf := NewLockFile(lockPath)

	// First lock should succeed
	err := lf.TryLock()
	if err != nil {
		t.Fatalf("First lock failed: %v", err)
	}

	// Verify PID file was created
	if _, err := os.Stat(lockPath); os.IsNotExist(err) {
		t.Fatal("Lock file was not created")
	}

	// Verify PID is correct
	pid, err := lf.readPID()
	if err != nil {
		t.Fatalf("Failed to read PID: %v", err)
	}
	if pid != os.Getpid() {
		t.Fatalf("Expected PID %d, got %d", os.Getpid(), pid)
	}

	// Second lock should fail (same process trying to lock again)
	lf2 := NewLockFile(lockPath)
	err = lf2.TryLock()
	if err == nil {
		t.Fatal("Second lock should have failed")
	}

	// Release lock
	err = lf.Release()
	if err != nil {
		t.Fatalf("Failed to release lock: %v", err)
	}

	// Verify lock file was removed
	if _, err := os.Stat(lockPath); !os.IsNotExist(err) {
		t.Fatal("Lock file was not removed")
	}

	// Third lock should succeed after release
	lf3 := NewLockFile(lockPath)
	err = lf3.TryLock()
	if err != nil {
		t.Fatalf("Lock after release failed: %v", err)
	}
	lf3.Release()
}

func TestLockFile_StaleLock(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "stale.pid")

	// Create a stale lock file with a non-existent PID
	stalePID := 999999
	err := os.WriteFile(lockPath, []byte("999999\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create stale lock: %v", err)
	}

	// Verify the PID doesn't exist
	lf := NewLockFile(lockPath)
	if lf.isProcessRunning(stalePID) {
		t.Skip("Test PID 999999 is actually running, skipping test")
	}

	// Lock should succeed and clean up stale lock
	err = lf.TryLock()
	if err != nil {
		t.Fatalf("Failed to acquire lock with stale lock file: %v", err)
	}

	// Verify new PID is written
	pid, err := lf.readPID()
	if err != nil {
		t.Fatalf("Failed to read PID: %v", err)
	}
	if pid != os.Getpid() {
		t.Fatalf("Expected PID %d, got %d", os.Getpid(), pid)
	}

	lf.Release()
}

func TestLockFile_CorruptedLock(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "corrupted.pid")

	// Create a corrupted lock file
	err := os.WriteFile(lockPath, []byte("not-a-number\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create corrupted lock: %v", err)
	}

	// Lock should succeed and clean up corrupted lock
	lf := NewLockFile(lockPath)
	err = lf.TryLock()
	if err != nil {
		t.Fatalf("Failed to acquire lock with corrupted lock file: %v", err)
	}

	// Verify new PID is written
	pid, err := lf.readPID()
	if err != nil {
		t.Fatalf("Failed to read PID: %v", err)
	}
	if pid != os.Getpid() {
		t.Fatalf("Expected PID %d, got %d", os.Getpid(), pid)
	}

	lf.Release()
}

func TestLockFile_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "concurrent.pid")

	lf1 := NewLockFile(lockPath)
	err := lf1.TryLock()
	if err != nil {
		t.Fatalf("First lock failed: %v", err)
	}
	defer lf1.Release()

	// Simulate concurrent access
	done := make(chan bool)
	go func() {
		lf2 := NewLockFile(lockPath)
		err := lf2.TryLock()
		if err == nil {
			t.Error("Concurrent lock should have failed")
		}
		done <- true
	}()

	select {
	case <-done:
		// Test passed
	case <-time.After(2 * time.Second):
		t.Fatal("Concurrent lock test timed out")
	}
}
