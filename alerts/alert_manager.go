package alerts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tomurashigaraki22/watchup-agent/detectors"
)

type ProcessPayload struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	CPU        float64 `json:"cpu"`
	Memory     float32 `json:"memory"`
}

type AlertPayload struct {
	ServerKey    string           `json:"server_key"`
	Metric       string           `json:"metric"`
	Usage        float64          `json:"usage"`
	Duration     int              `json:"duration"`
	TopProcesses []ProcessPayload `json:"top_processes"`
	Timestamp    string           `json:"timestamp"`
}

type AlertManager struct {
	serverKey string
	onAlert   func(AlertPayload)
}

func NewAlertManager(serverKey string, onAlert func(AlertPayload)) *AlertManager {
	return &AlertManager{
		serverKey: serverKey,
		onAlert:   onAlert,
	}
}

func (am *AlertManager) HandleSpikeEvent(event detectors.SpikeEvent) {
	processes := make([]ProcessPayload, 0, len(event.TopProcesses))
	for _, proc := range event.TopProcesses {
		processes = append(processes, ProcessPayload{
			PID:    proc.PID,
			Name:   proc.Name,
			CPU:    proc.CPUPercent,
			Memory: proc.MemPercent,
		})
	}

	alert := AlertPayload{
		ServerKey:    am.serverKey,
		Metric:       string(event.Metric),
		Usage:        event.Value,
		Duration:     int(event.Duration.Seconds()),
		TopProcesses: processes,
		Timestamp:    event.Timestamp.Format(time.RFC3339),
	}

	// Log alert
	alertJSON, _ := json.MarshalIndent(alert, "", "  ")
	fmt.Printf("\n🚨 ALERT TRIGGERED:\n%s\n\n", string(alertJSON))

	// Call handler if provided
	if am.onAlert != nil {
		am.onAlert(alert)
	}
}
