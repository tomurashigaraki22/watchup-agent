package collectors

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUMetrics struct {
	UsagePercent float64
	PerCore      []float64
}

type CPUCollector struct{}

func NewCPUCollector() *CPUCollector {
	return &CPUCollector{}
}

func (c *CPUCollector) Collect() (*CPUMetrics, error) {
	// Get total CPU usage
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU usage: %w", err)
	}

	if len(percentages) == 0 {
		return nil, fmt.Errorf("no CPU usage data available")
	}

	// Get per-core usage (optional)
	perCore, err := cpu.Percent(0, true)
	if err != nil {
		perCore = []float64{}
	}

	return &CPUMetrics{
		UsagePercent: percentages[0],
		PerCore:      perCore,
	}, nil
}
