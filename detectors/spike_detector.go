package detectors

import (
	"fmt"
	"time"

	"github.com/tomurashigaraki22/watchup-agent/collectors"
)

type MetricType string

const (
	MetricCPU        MetricType = "cpu"
	MetricRAM        MetricType = "ram"
	MetricProcessCPU MetricType = "process_cpu"
)

type ThresholdConfig struct {
	Threshold      float64
	Duration       time.Duration
	RequiredSamples int
}

type SpikeEvent struct {
	Metric       MetricType
	Value        float64
	Threshold    float64
	Duration     time.Duration
	TopProcesses []collectors.ProcessInfo
	Timestamp    time.Time
}

type metricState struct {
	violationCount int
	lastValue      float64
}

type SpikeDetector struct {
	cpuConfig     ThresholdConfig
	ramConfig     ThresholdConfig
	processCPUConfig ThresholdConfig
	samplingInterval time.Duration
	
	cpuState     metricState
	ramState     metricState
	processCPUState metricState
	
	onAlert func(SpikeEvent)
}

func NewSpikeDetector(
	cpuThreshold, ramThreshold, processCPUThreshold float64,
	cpuDuration, ramDuration, processCPUDuration time.Duration,
	samplingInterval time.Duration,
	onAlert func(SpikeEvent),
) *SpikeDetector {
	return &SpikeDetector{
		cpuConfig: ThresholdConfig{
			Threshold:      cpuThreshold,
			Duration:       cpuDuration,
			RequiredSamples: int(cpuDuration / samplingInterval),
		},
		ramConfig: ThresholdConfig{
			Threshold:      ramThreshold,
			Duration:       ramDuration,
			RequiredSamples: int(ramDuration / samplingInterval),
		},
		processCPUConfig: ThresholdConfig{
			Threshold:      processCPUThreshold,
			Duration:       processCPUDuration,
			RequiredSamples: int(processCPUDuration / samplingInterval),
		},
		samplingInterval: samplingInterval,
		onAlert:          onAlert,
	}
}

func (sd *SpikeDetector) CheckCPU(usage float64, processes []collectors.ProcessInfo) {
	sd.cpuState.lastValue = usage
	
	if usage > sd.cpuConfig.Threshold {
		sd.cpuState.violationCount++
		
		if sd.cpuState.violationCount >= sd.cpuConfig.RequiredSamples {
			if sd.onAlert != nil {
				sd.onAlert(SpikeEvent{
					Metric:       MetricCPU,
					Value:        usage,
					Threshold:    sd.cpuConfig.Threshold,
					Duration:     sd.cpuConfig.Duration,
					TopProcesses: processes,
					Timestamp:    time.Now(),
				})
			}
			// Reset after alert to avoid duplicate alerts
			sd.cpuState.violationCount = 0
		}
	} else {
		sd.cpuState.violationCount = 0
	}
}

func (sd *SpikeDetector) CheckRAM(usage float64, processes []collectors.ProcessInfo) {
	sd.ramState.lastValue = usage
	
	if usage > sd.ramConfig.Threshold {
		sd.ramState.violationCount++
		
		if sd.ramState.violationCount >= sd.ramConfig.RequiredSamples {
			if sd.onAlert != nil {
				sd.onAlert(SpikeEvent{
					Metric:       MetricRAM,
					Value:        usage,
					Threshold:    sd.ramConfig.Threshold,
					Duration:     sd.ramConfig.Duration,
					TopProcesses: processes,
					Timestamp:    time.Now(),
				})
			}
			sd.ramState.violationCount = 0
		}
	} else {
		sd.ramState.violationCount = 0
	}
}

func (sd *SpikeDetector) CheckProcessCPU(processes []collectors.ProcessInfo) {
	if len(processes) == 0 {
		sd.processCPUState.violationCount = 0
		return
	}
	
	topProcess := processes[0]
	sd.processCPUState.lastValue = topProcess.CPUPercent
	
	if topProcess.CPUPercent > sd.processCPUConfig.Threshold {
		sd.processCPUState.violationCount++
		
		if sd.processCPUState.violationCount >= sd.processCPUConfig.RequiredSamples {
			if sd.onAlert != nil {
				sd.onAlert(SpikeEvent{
					Metric:       MetricProcessCPU,
					Value:        topProcess.CPUPercent,
					Threshold:    sd.processCPUConfig.Threshold,
					Duration:     sd.processCPUConfig.Duration,
					TopProcesses: processes,
					Timestamp:    time.Now(),
				})
			}
			sd.processCPUState.violationCount = 0
		}
	} else {
		sd.processCPUState.violationCount = 0
	}
}

func (sd *SpikeDetector) GetStatus() string {
	return fmt.Sprintf(
		"CPU: %.1f%% (%d/%d samples), RAM: %.1f%% (%d/%d samples), Process: %.1f%% (%d/%d samples)",
		sd.cpuState.lastValue, sd.cpuState.violationCount, sd.cpuConfig.RequiredSamples,
		sd.ramState.lastValue, sd.ramState.violationCount, sd.ramConfig.RequiredSamples,
		sd.processCPUState.lastValue, sd.processCPUState.violationCount, sd.processCPUConfig.RequiredSamples,
	)
}

// UpdateThresholds updates the detector's threshold configuration dynamically
func (sd *SpikeDetector) UpdateThresholds(
	cpuThreshold, ramThreshold, processCPUThreshold float64,
	cpuDuration, ramDuration, processCPUDuration time.Duration,
) {
	sd.cpuConfig.Threshold = cpuThreshold
	sd.cpuConfig.Duration = cpuDuration
	sd.cpuConfig.RequiredSamples = int(cpuDuration / sd.samplingInterval)
	
	sd.ramConfig.Threshold = ramThreshold
	sd.ramConfig.Duration = ramDuration
	sd.ramConfig.RequiredSamples = int(ramDuration / sd.samplingInterval)
	
	sd.processCPUConfig.Threshold = processCPUThreshold
	sd.processCPUConfig.Duration = processCPUDuration
	sd.processCPUConfig.RequiredSamples = int(processCPUDuration / sd.samplingInterval)
	
	// Reset violation counts when thresholds change
	sd.cpuState.violationCount = 0
	sd.ramState.violationCount = 0
	sd.processCPUState.violationCount = 0
}
