package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const (
	DefaultLockFile = "/var/run/watchup-agent.pid"
	FallbackLockDir = "/tmp"
)

type LockFile struct {
	path string
}

// NewLockFile creates a new lock file manager
func NewLockFile(path string) *LockFile {
	if path == "" {
		path = DefaultLockFile
	}
	return &LockFile{path: path}
}

// TryLock attempts to acquire the lock
// Returns error if another instance is running
func (lf *LockFile) TryLock() error {
	// Check if lock file exists
	if _, err := os.Stat(lf.path); err == nil {
		// Lock file exists, check if process is still running
		pid, err := lf.readPID()
		if err != nil {
			// Corrupted lock file, clean it up
			fmt.Printf("Warning: Corrupted lock file found, cleaning up...\n")
			os.Remove(lf.path)
		} else if lf.isProcessRunning(pid) {
			return fmt.Errorf("another instance of watchup-agent is already running (PID: %d)", pid)
		} else {
			// Stale lock file, clean it up
			fmt.Printf("Warning: Stale lock file found (PID: %d no longer running), cleaning up...\n", pid)
			os.Remove(lf.path)
		}
	}

	// Try to create lock file
	if err := lf.writePID(); err != nil {
		// If we can't write to default location, try fallback
		if lf.path == DefaultLockFile {
			fallbackPath := filepath.Join(FallbackLockDir, "watchup-agent.pid")
			fmt.Printf("Warning: Cannot write to %s, using fallback: %s\n", DefaultLockFile, fallbackPath)
			lf.path = fallbackPath
			return lf.writePID()
		}
		return fmt.Errorf("failed to create lock file: %w", err)
	}

	return nil
}

// Release removes the lock file
func (lf *LockFile) Release() error {
	if err := os.Remove(lf.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove lock file: %w", err)
	}
	return nil
}

// readPID reads the PID from the lock file
func (lf *LockFile) readPID() (int, error) {
	data, err := os.ReadFile(lf.path)
	if err != nil {
		return 0, err
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("invalid PID in lock file: %s", pidStr)
	}

	return pid, nil
}

// writePID writes the current process PID to the lock file
func (lf *LockFile) writePID() error {
	// Ensure directory exists
	dir := filepath.Dir(lf.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	pid := os.Getpid()
	return os.WriteFile(lf.path, []byte(fmt.Sprintf("%d\n", pid)), 0644)
}

// isProcessRunning checks if a process with the given PID is running
func (lf *LockFile) isProcessRunning(pid int) bool {
	// Try to send signal 0 to check if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix systems, signal 0 checks if process exists without actually sending a signal
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// GetPath returns the lock file path
func (lf *LockFile) GetPath() string {
	return lf.path
}
