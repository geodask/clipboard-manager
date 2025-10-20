package daemon

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type PIDFile struct {
	path string
}

func NewPIDFile(path string) *PIDFile {
	return &PIDFile{path: path}
}

func (p *PIDFile) Create() error {
	running, oldPID, err := p.IsRunning()

	if err != nil {
		return fmt.Errorf("failed to check existing PID file: %w", err)
	}

	if running {
		return fmt.Errorf("daemon already running with PID %d", oldPID)
	}

	if oldPID != 0 {
		os.Remove(p.path)
	}

	pid := os.Getpid()
	content := fmt.Sprintf("%d\n", pid)

	if err := os.WriteFile(p.path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	return nil
}

func (p *PIDFile) Remove() error {
	if err := os.Remove(p.path); err != nil && os.IsExist(err) {
		return fmt.Errorf("failed to remove PID file: %w", err)
	}

	return nil
}

func (p *PIDFile) Read() (int, error) {
	data, err := os.ReadFile(p.path)
	if os.IsNotExist(err) {
		return 0, nil
	}

	if err != nil {
		return 0, fmt.Errorf("failed to read PID file: %w", err)
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("invalid PID in file: %w", err)
	}

	return pid, nil

}

func (p *PIDFile) IsRunning() (bool, int, error) {
	pid, err := p.Read()
	if err != nil {
		return false, 0, err
	}

	if pid == 0 {
		return false, 0, nil
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false, pid, nil
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false, pid, nil
	}

	return true, pid, nil

}
