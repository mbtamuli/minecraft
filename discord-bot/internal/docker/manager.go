package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// Manager handles Docker Compose operations
type Manager struct {
	composePath string
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewManager creates a new Docker manager instance
func NewManager(composePath string) *Manager {
	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		composePath: composePath,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// runDockerCompose executes a docker compose command and returns the output
func (m *Manager) runDockerCompose(args ...string) ([]byte, error) {
	// Build the command - Docker Compose will automatically find compose.yaml in the working directory
	cmdArgs := []string{"compose"}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.CommandContext(m.ctx, "docker", cmdArgs...)
	cmd.Dir = filepath.Dir(m.composePath)

	output, err := cmd.Output()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("docker compose command failed: %s, stderr: %s", err, string(exitError.Stderr))
		}
		return nil, fmt.Errorf("failed to execute docker compose command: %w", err)
	}

	return output, nil
}

// ServiceStatus represents the status of a service
type ServiceStatus struct {
	IsRunning bool
	Uptime    string
}

// GetStatus returns the current status of Docker services
func (m *Manager) GetStatus() (*ServiceStatus, error) {
	output, err := m.runDockerCompose("ps", "--format", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to get service status: %w", err)
	}

	// Parse JSON output if needed
	var services []map[string]any
	if len(output) > 0 {
		// Split by lines since docker compose ps --format json outputs one JSON object per line
		lines := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")
		for line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var service map[string]any
			if err := json.Unmarshal([]byte(line), &service); err != nil {
				return nil, fmt.Errorf("failed to parse service status JSON: %w", err)
			}
			services = append(services, service)
		}
	}

	// Look for the "mc" service
	for _, service := range services {
		if serviceName, ok := service["Service"].(string); ok && serviceName == "mc" {
			status := &ServiceStatus{}

			// Check if service is running
			if state, ok := service["State"].(string); ok {
				status.IsRunning = state == "running"
			}

			// Get uptime information
			if runningFor, ok := service["RunningFor"].(string); ok {
				status.Uptime = runningFor
			}

			return status, nil
		}
	}

	// If mc service not found, return not running
	return &ServiceStatus{IsRunning: false, Uptime: ""}, nil
}

// StartServices starts all Docker services
func (m *Manager) StartServices() error {
	log.Println("Starting Docker services...")
	_, err := m.runDockerCompose("up", "--wait")
	if err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}
	log.Println("Docker services started successfully")
	return nil
}

// StopServices stops all Docker services
func (m *Manager) StopServices() error {
	log.Println("Stopping Docker services...")
	_, err := m.runDockerCompose("down")
	if err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}
	log.Println("Docker services stopped successfully")
	return nil
}

// Close cleans up resources and cancels any ongoing operations
func (m *Manager) Close() error {
	if m.cancel != nil {
		m.cancel()
	}
	// TODO: Add any additional cleanup logic
	return nil
}
