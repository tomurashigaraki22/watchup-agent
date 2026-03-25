package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tomurashigaraki22/watchup-agent/alerts"
	"github.com/tomurashigaraki22/watchup-agent/collectors"
	"github.com/tomurashigaraki22/watchup-agent/config"
	"github.com/tomurashigaraki22/watchup-agent/detectors"
	"github.com/tomurashigaraki22/watchup-agent/internal"
	"github.com/tomurashigaraki22/watchup-agent/transport"
)

const defaultConfigPath = "/etc/watchup/config.yaml"

func main() {
	fmt.Println("Watchup Server Agent started")

	// Load configuration
	configPath := defaultConfigPath
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Check if agent is registered
	registrar := internal.NewRegistrar(cfg.APIEndpoint)
	if !registrar.CheckRegistration(cfg) {
		fmt.Println("\n⚠️  Agent is not registered.")
		
		if err := registrar.PerformRegistration(cfg, configPath); err != nil {
			fmt.Printf("\n❌ Registration failed: %v\n", err)
			fmt.Println("\nThe agent cannot start without registration.")
			fmt.Println("Please check your credentials and try again.")
			os.Exit(1)
		}
	}

	fmt.Printf("\n✓ Agent registered successfully\n")
	fmt.Printf("Project ID: %s\n", cfg.ProjectID)
	fmt.Printf("Server: %s\n", cfg.ServerIdentifier)
	fmt.Printf("Sampling interval: %ds\n", cfg.SamplingInterval)
	fmt.Printf("Alert Thresholds - CPU: %d%%, RAM: %d%%\n\n", 
		cfg.Alerts.CPU.Threshold, cfg.Alerts.RAM.Threshold)

	// Initialize collectors
	cpuCollector := collectors.NewCPUCollector()
	memCollector := collectors.NewMemoryCollector()
	procCollector := collectors.NewProcessCollector()

	// Initialize API client
	apiClient := transport.NewAPIClient(cfg)

	// Initialize alert manager
	alertManager := alerts.NewAlertManager(cfg.ServerKey, func(alert alerts.AlertPayload) {
		if err := apiClient.SendAlert(alert); err != nil {
			fmt.Printf("Failed to send alert: %v\n", err)
		}
	})

	// Initialize spike detector
	spikeDetector := detectors.NewSpikeDetector(
		float64(cfg.Alerts.CPU.Threshold),
		float64(cfg.Alerts.RAM.Threshold),
		float64(cfg.Alerts.ProcessCPU.Threshold),
		time.Duration(cfg.Alerts.CPU.Duration)*time.Second,
		time.Duration(cfg.Alerts.RAM.Duration)*time.Second,
		time.Duration(cfg.Alerts.ProcessCPU.Duration)*time.Second,
		cfg.GetSamplingDuration(),
		alertManager.HandleSpikeEvent,
	)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutdown signal received")
		cancel()
	}()

	// Monitoring tick function
	tickCount := 0
	monitoringTick := func() error {
		tickCount++

		// Collect CPU metrics
		cpuMetrics, err := cpuCollector.Collect()
		if err != nil {
			return fmt.Errorf("CPU collection failed: %w", err)
		}

		// Collect memory metrics
		memMetrics, err := memCollector.Collect()
		if err != nil {
			return fmt.Errorf("memory collection failed: %w", err)
		}

		// Collect top processes
		topProcesses, err := procCollector.CollectTopCPU(5)
		if err != nil {
			return fmt.Errorf("process collection failed: %w", err)
		}

		// Check for spikes
		spikeDetector.CheckCPU(cpuMetrics.UsagePercent, topProcesses)
		spikeDetector.CheckRAM(memMetrics.UsedPercent, topProcesses)
		spikeDetector.CheckProcessCPU(topProcesses)

		// Log status every 12 ticks (60 seconds at 5s interval)
		if tickCount%12 == 0 {
			fmt.Printf("[%s] CPU: %.1f%%, RAM: %.1f%% | %s\n",
				time.Now().Format("15:04:05"),
				cpuMetrics.UsagePercent,
				memMetrics.UsedPercent,
				spikeDetector.GetStatus(),
			)
		}

		// Send metrics to API (every tick)
		if err := apiClient.SendMetrics(cpuMetrics.UsagePercent, memMetrics.UsedPercent); err != nil {
			fmt.Printf("Failed to send metrics: %v\n", err)
		}

		return nil
	}

	// Config reload function
	configReload := func() error {
		fmt.Println("Checking for configuration updates...")
		newCfg, err := apiClient.FetchConfig()
		if err != nil {
			fmt.Printf("Failed to fetch config: %v\n", err)
			return nil // Don't fail on config fetch errors
		}

		// Update configuration (simplified - in production, would need proper synchronization)
		*cfg = *newCfg
		fmt.Println("Configuration updated from API")
		return nil
	}

	// Initialize scheduler
	scheduler := internal.NewScheduler(cfg.GetSamplingDuration(), monitoringTick)
	scheduler.SetConfigReloadHandler(configReload)

	// Start monitoring loop
	fmt.Println("Starting monitoring loop...")
	if err := scheduler.Start(ctx); err != nil && err != context.Canceled {
		fmt.Printf("Scheduler error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Watchup Server Agent stopped")
}
