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
	fmt.Println("Watchup Server Agent starting...")

	// Acquire instance lock to prevent multiple instances
	lockFile := internal.NewLockFile("")
	if err := lockFile.TryLock(); err != nil {
		fmt.Printf("\n❌ %v\n", err)
		fmt.Println("\nOnly one instance of watchup-agent can run at a time.")
		fmt.Println("If you believe this is an error, check for stale processes:")
		fmt.Println("  ps aux | grep watchup-agent")
		os.Exit(1)
	}
	defer lockFile.Release()

	fmt.Printf("✓ Instance lock acquired (PID: %d)\n", os.Getpid())

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
	fmt.Printf("API Endpoint: %s\n", cfg.APIEndpoint)
	fmt.Printf("\n--- Monitoring Configuration ---\n")
	fmt.Printf("Sampling interval: %ds (metrics collected every %d seconds)\n", 
		cfg.SamplingInterval, cfg.SamplingInterval)
	fmt.Printf("Config reload: every 60s\n")
	fmt.Printf("\n--- Alert Thresholds ---\n")
	fmt.Printf("CPU: %d%% sustained for %ds (%d samples)\n", 
		cfg.Alerts.CPU.Threshold, 
		cfg.Alerts.CPU.Duration,
		cfg.Alerts.CPU.Duration/cfg.SamplingInterval)
	fmt.Printf("RAM: %d%% sustained for %ds (%d samples)\n", 
		cfg.Alerts.RAM.Threshold,
		cfg.Alerts.RAM.Duration,
		cfg.Alerts.RAM.Duration/cfg.SamplingInterval)
	fmt.Printf("Process CPU: %d%% sustained for %ds (%d samples)\n\n",
		cfg.Alerts.ProcessCPU.Threshold,
		cfg.Alerts.ProcessCPU.Duration,
		cfg.Alerts.ProcessCPU.Duration/cfg.SamplingInterval)

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
		fmt.Println("\nShutdown signal received, cleaning up...")
		lockFile.Release()
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
			fmt.Printf("[MONITOR] [%s] CPU: %.1f%%, RAM: %.1f%% | %s\n",
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
		fmt.Printf("[CONFIG] Checking for configuration updates...\n")
		newCfg, err := apiClient.FetchConfig()
		if err != nil {
			fmt.Printf("[CONFIG] Failed to fetch config: %v\n", err)
			return nil // Don't fail on config fetch errors
		}

		// Only update alert thresholds, preserve local settings
		// (sampling_interval, server_key, project_id, etc. should not be changed by API)
		oldSamplingInterval := cfg.SamplingInterval
		oldServerKey := cfg.ServerKey
		oldProjectID := cfg.ProjectID
		oldServerIdentifier := cfg.ServerIdentifier
		oldAPIEndpoint := cfg.APIEndpoint
		oldRegistered := cfg.Registered

		// Warn if API is trying to change sampling interval
		if newCfg.SamplingInterval != 0 && newCfg.SamplingInterval != oldSamplingInterval {
			fmt.Printf("[CONFIG] Warning: API returned sampling_interval=%ds, keeping local value=%ds\n",
				newCfg.SamplingInterval, oldSamplingInterval)
		}

		// Update only the alerts configuration
		cfg.Alerts = newCfg.Alerts

		// Restore local settings
		cfg.SamplingInterval = oldSamplingInterval
		cfg.ServerKey = oldServerKey
		cfg.ProjectID = oldProjectID
		cfg.ServerIdentifier = oldServerIdentifier
		cfg.APIEndpoint = oldAPIEndpoint
		cfg.Registered = oldRegistered

		// Calculate sample counts for logging
		cpuSamples := cfg.Alerts.CPU.Duration / cfg.SamplingInterval
		ramSamples := cfg.Alerts.RAM.Duration / cfg.SamplingInterval
		processSamples := cfg.Alerts.ProcessCPU.Duration / cfg.SamplingInterval

		fmt.Printf("[CONFIG] Alert thresholds updated:\n")
		fmt.Printf("  CPU: %d%% for %ds (%d samples)\n", 
			cfg.Alerts.CPU.Threshold, cfg.Alerts.CPU.Duration, cpuSamples)
		fmt.Printf("  RAM: %d%% for %ds (%d samples)\n", 
			cfg.Alerts.RAM.Threshold, cfg.Alerts.RAM.Duration, ramSamples)
		fmt.Printf("  Process: %d%% for %ds (%d samples)\n",
			cfg.Alerts.ProcessCPU.Threshold, cfg.Alerts.ProcessCPU.Duration, processSamples)

		// Update spike detector with new thresholds
		spikeDetector.UpdateThresholds(
			float64(cfg.Alerts.CPU.Threshold),
			float64(cfg.Alerts.RAM.Threshold),
			float64(cfg.Alerts.ProcessCPU.Threshold),
			time.Duration(cfg.Alerts.CPU.Duration)*time.Second,
			time.Duration(cfg.Alerts.RAM.Duration)*time.Second,
			time.Duration(cfg.Alerts.ProcessCPU.Duration)*time.Second,
		)

		// Immediately collect and check metrics after config update
		cpuMetrics, _ := cpuCollector.Collect()
		memMetrics, _ := memCollector.Collect()
		topProcesses, _ := procCollector.CollectTopCPU(5)
		
		if cpuMetrics != nil && memMetrics != nil && topProcesses != nil {
			spikeDetector.CheckCPU(cpuMetrics.UsagePercent, topProcesses)
			spikeDetector.CheckRAM(memMetrics.UsedPercent, topProcesses)
			spikeDetector.CheckProcessCPU(topProcesses)
			fmt.Printf("[CONFIG] Resource check: CPU: %.1f%%, RAM: %.1f%%\n", 
				cpuMetrics.UsagePercent, memMetrics.UsedPercent)
		}

		return nil
	}

	// Initialize scheduler
	scheduler := internal.NewScheduler(cfg.GetSamplingDuration(), monitoringTick)
	scheduler.SetConfigReloadHandler(configReload)

	// Start monitoring loop
	fmt.Printf("Starting monitoring loop (interval: %v)...\n", cfg.GetSamplingDuration())
	if err := scheduler.Start(ctx); err != nil && err != context.Canceled {
		fmt.Printf("Scheduler error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Watchup Server Agent stopped")
}
