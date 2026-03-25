package collectors

import (
	"fmt"
	"sort"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcessInfo struct {
	PID         int32
	Name        string
	CPUPercent  float64
	MemPercent  float32
}

type ProcessCollector struct{}

func NewProcessCollector() *ProcessCollector {
	return &ProcessCollector{}
}

func (p *ProcessCollector) Collect(topN int) ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var processInfos []ProcessInfo

	for _, proc := range processes {
		name, err := proc.Name()
		if err != nil {
			continue
		}

		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}

		memPercent, err := proc.MemoryPercent()
		if err != nil {
			memPercent = 0
		}

		processInfos = append(processInfos, ProcessInfo{
			PID:        proc.Pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
		})
	}

	// Sort by CPU usage (descending)
	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].CPUPercent > processInfos[j].CPUPercent
	})

	// Return top N processes
	if len(processInfos) > topN {
		processInfos = processInfos[:topN]
	}

	return processInfos, nil
}

func (p *ProcessCollector) CollectTopCPU(topN int) ([]ProcessInfo, error) {
	return p.Collect(topN)
}

func (p *ProcessCollector) CollectTopMemory(topN int) ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	var processInfos []ProcessInfo

	for _, proc := range processes {
		name, err := proc.Name()
		if err != nil {
			continue
		}

		cpuPercent, err := proc.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}

		memPercent, err := proc.MemoryPercent()
		if err != nil {
			memPercent = 0
		}

		processInfos = append(processInfos, ProcessInfo{
			PID:        proc.Pid,
			Name:       name,
			CPUPercent: cpuPercent,
			MemPercent: memPercent,
		})
	}

	// Sort by memory usage (descending)
	sort.Slice(processInfos, func(i, j int) bool {
		return processInfos[i].MemPercent > processInfos[j].MemPercent
	})

	// Return top N processes
	if len(processInfos) > topN {
		processInfos = processInfos[:topN]
	}

	return processInfos, nil
}
