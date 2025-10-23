package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// Check if clipboard daemon is running
func isDaemonRunning() bool {
	// Check for running clipboard-manager processes
	cmd := exec.Command("pgrep", "-f", "clipboard-manager daemon")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	pids := strings.TrimSpace(string(output))
	return pids != ""
}

// Start daemon if not running
func ensureDaemonRunning() {
	if isDaemonRunning() {
		fmt.Println("✓ Clipboard daemon is already running")
		return
	}
	
	fmt.Println("Starting clipboard daemon...")
	
	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		execPath = "./clipboard-manager"
	}
	
	// Start daemon in background
	cmd := exec.Command(execPath, "daemon")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true, // Create new process group
	}
	
	// Redirect output to avoid blocking
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil
	
	if err := cmd.Start(); err != nil {
		fmt.Printf("Warning: Could not start daemon: %v\n", err)
		return
	}
	
	// Don't wait for the process, let it run independently
	go func() {
		cmd.Wait()
	}()
	
	fmt.Printf("✓ Clipboard daemon started (PID: %d)\n", cmd.Process.Pid)
	
	// Save PID for later reference
	saveDaemonPID(cmd.Process.Pid)
}

// Save daemon PID to file
func saveDaemonPID(pid int) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	
	pidDir := filepath.Join(home, ".local", "share", "clipboard-manager")
	os.MkdirAll(pidDir, 0755)
	
	pidFile := filepath.Join(pidDir, "daemon.pid")
	os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

// Get daemon PID from file
func getDaemonPID() int {
	home, err := os.UserHomeDir()
	if err != nil {
		return 0
	}
	
	pidFile := filepath.Join(home, ".local", "share", "clipboard-manager", "daemon.pid")
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0
	}
	
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	
	return pid
}

// Stop daemon
func stopDaemon() {
	pid := getDaemonPID()
	if pid == 0 {
		fmt.Println("No daemon PID found")
		return
	}
	
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("Could not find process %d: %v\n", pid, err)
		return
	}
	
	if err := process.Signal(syscall.SIGTERM); err != nil {
		fmt.Printf("Could not stop daemon: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Daemon stopped (PID: %d)\n", pid)
	
	// Remove PID file
	home, _ := os.UserHomeDir()
	pidFile := filepath.Join(home, ".local", "share", "clipboard-manager", "daemon.pid")
	os.Remove(pidFile)
}

// Show daemon status
func showDaemonStatus() {
	if isDaemonRunning() {
		pid := getDaemonPID()
		fmt.Printf("✓ Clipboard daemon is running (PID: %d)\n", pid)
	} else {
		fmt.Println("✗ Clipboard daemon is not running")
		fmt.Println("Start it with: ./clipboard-manager daemon")
	}
}