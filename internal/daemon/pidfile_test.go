// internal/daemon/pidfile_test.go
package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestPIDFile_Create(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(path string)
		wantErr bool
	}{
		{
			name:    "Success - No existing file",
			setup:   func(path string) {},
			wantErr: false,
		},
		{
			name: "Success - Stale file cleanup",
			setup: func(path string) {
				os.WriteFile(path, []byte("99999\n"), 0644)
			},
			wantErr: false,
		},
		{
			name: "Error - Another instance running",
			setup: func(path string) {
				pid := os.Getpid()
				os.WriteFile(path, []byte(fmt.Sprintf("%d\n", pid)), 0644)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			pidPath := filepath.Join(tmpDir, "test.pid")

			if tt.setup != nil {
				tt.setup(pidPath)
			}

			pidFile := NewPIDFile(pidPath)
			err := pidFile.Create()

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if _, err := os.Stat(pidPath); os.IsNotExist(err) {
					t.Error("PID file was not created")
				}

				pid, err := pidFile.Read()
				if err != nil {
					t.Errorf("failed to read PID: %v", err)
				}
				if pid != os.Getpid() {
					t.Errorf("PID = %d, want %d", pid, os.Getpid())
				}

				pidFile.Remove()
			}
		})
	}
}

func TestPIDFile_Remove(t *testing.T) {
	tmpDir := t.TempDir()
	pidPath := filepath.Join(tmpDir, "test.pid")

	pidFile := NewPIDFile(pidPath)

	if err := pidFile.Create(); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	if _, err := os.Stat(pidPath); os.IsNotExist(err) {
		t.Fatal("PID file was not created")
	}

	if err := pidFile.Remove(); err != nil {
		t.Errorf("Remove() error = %v", err)
	}

	if _, err := os.Stat(pidPath); !os.IsNotExist(err) {
		t.Error("PID file still exists after Remove()")
	}
}

func TestPIDFile_Read(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantPID int
		wantErr bool
	}{
		{
			name:    "Valid PID",
			content: "12345\n",
			wantPID: 12345,
			wantErr: false,
		},
		{
			name:    "Valid PID without newline",
			content: "67890",
			wantPID: 67890,
			wantErr: false,
		},
		{
			name:    "Invalid PID",
			content: "not-a-number\n",
			wantPID: 0,
			wantErr: true,
		},
		{
			name:    "Empty file",
			content: "",
			wantPID: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			pidPath := filepath.Join(tmpDir, "test.pid")

			os.WriteFile(pidPath, []byte(tt.content), 0644)

			pidFile := NewPIDFile(pidPath)
			pid, err := pidFile.Read()

			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if pid != tt.wantPID {
				t.Errorf("Read() = %d, want %d", pid, tt.wantPID)
			}
		})
	}
}

func TestPIDFile_IsRunning(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(path string)
		wantRunning bool
		wantPID     int
	}{
		{
			name:        "No PID file",
			setup:       func(path string) {},
			wantRunning: false,
			wantPID:     0,
		},
		{
			name: "Current process (running)",
			setup: func(path string) {
				pid := os.Getpid()
				os.WriteFile(path, []byte(fmt.Sprintf("%d\n", pid)), 0644)
			},
			wantRunning: true,
			wantPID:     os.Getpid(),
		},
		{
			name: "Non-existent PID (stale)",
			setup: func(path string) {
				os.WriteFile(path, []byte("99999\n"), 0644)
			},
			wantRunning: false,
			wantPID:     99999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			pidPath := filepath.Join(tmpDir, "test.pid")

			if tt.setup != nil {
				tt.setup(pidPath)
			}

			pidFile := NewPIDFile(pidPath)
			running, pid, err := pidFile.IsRunning()

			if err != nil {
				t.Errorf("IsRunning() error = %v", err)
				return
			}

			if running != tt.wantRunning {
				t.Errorf("IsRunning() running = %v, want %v", running, tt.wantRunning)
			}

			if pid != tt.wantPID {
				t.Errorf("IsRunning() pid = %d, want %d", pid, tt.wantPID)
			}
		})
	}
}

func TestPIDFile_CreateRemoveCycle(t *testing.T) {
	tmpDir := t.TempDir()
	pidPath := filepath.Join(tmpDir, "test.pid")

	pidFile := NewPIDFile(pidPath)

	if err := pidFile.Create(); err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	running, pid, err := pidFile.IsRunning()
	if err != nil {
		t.Fatalf("IsRunning() error: %v", err)
	}
	if !running {
		t.Error("IsRunning() = false, want true")
	}
	if pid != os.Getpid() {
		t.Errorf("PID = %d, want %d", pid, os.Getpid())
	}

	if err := pidFile.Create(); err == nil {
		t.Error("Create() succeeded on second call, should fail")
	}

	if err := pidFile.Remove(); err != nil {
		t.Errorf("Remove() error = %v", err)
	}

	running, _, _ = pidFile.IsRunning()
	if running {
		t.Error("IsRunning() = true after Remove(), want false")
	}

	if err := pidFile.Create(); err != nil {
		t.Errorf("Create() after Remove() failed: %v", err)
	}

	pidFile.Remove()
}
