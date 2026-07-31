package worker

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"

	"github.com/google/uuid"
)

// WorkerCapabilities represents hardware specs and system toolchains detected on a worker node.
type WorkerCapabilities struct {
	WorkerID   string   `json:"worker_id"`
	Hostname   string   `json:"hostname"`
	IPAddress  string   `json:"ip_address"`
	OS         string   `json:"os"`
	Arch       string   `json:"arch"`
	CPUCores   int      `json:"cpu_cores"`
	Toolchains []string `json:"toolchains"`
}

// DetectCapabilities automatically probes system hardware and available CLI toolchains.
func DetectCapabilities() (*WorkerCapabilities, error) {
	workerUUID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID v7 for worker: %w", err)
	}

	hostname, _ := os.Hostname()
	ipAddress := getOutboundIP()

	caps := &WorkerCapabilities{
		WorkerID:   workerUUID.String(),
		Hostname:   hostname,
		IPAddress:  ipAddress,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		CPUCores:   runtime.NumCPU(),
		Toolchains: []string{},
	}

	// List of binaries to probe in system PATH
	toolsToProbe := []string{"mv", "cp", "rm", "ffmpeg", "ffprobe", "imagemagick", "convert", "python3", "7z", "docker"}

	for _, tool := range toolsToProbe {
		if _, err := exec.LookPath(tool); err == nil {
			caps.Toolchains = append(caps.Toolchains, tool)
		}
	}

	return caps, nil
}

// getOutboundIP finds the local IP address used for network outbound connections.
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
