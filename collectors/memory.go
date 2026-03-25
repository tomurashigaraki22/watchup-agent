package collectors

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryMetrics struct {
	Total        uint64
	Used         uint64
	Available    uint64
	UsedPercent  float64
}

type MemoryCollector struct{}

func NewMemoryCollector() *MemoryCollector {
	return &MemoryCollector{}
}

func (m *MemoryCollector) Collect() (*MemoryMetrics, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory stats: %w", err)
	}

	return &MemoryMetrics{
		Total:       vmStat.Total,
		Used:        vmStat.Used,
		Available:   vmStat.Available,
		UsedPercent: vmStat.UsedPercent,
	}, nil
}
