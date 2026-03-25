package internal

import (
	"context"
	"fmt"
	"time"
)

type Scheduler struct {
	interval       time.Duration
	configInterval time.Duration
	onTick         func() error
	onConfigReload func() error
}

func NewScheduler(interval time.Duration, onTick func() error) *Scheduler {
	return &Scheduler{
		interval:       interval,
		configInterval: 60 * time.Second,
		onTick:         onTick,
	}
}

func (s *Scheduler) SetConfigReloadHandler(handler func() error) {
	s.onConfigReload = handler
}

func (s *Scheduler) Start(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	configTicker := time.NewTicker(s.configInterval)
	defer configTicker.Stop()

	fmt.Printf("Scheduler started - Monitoring: every %v, Config reload: every %v\n", 
		s.interval, s.configInterval)

	// Run initial tick immediately
	if err := s.onTick(); err != nil {
		fmt.Printf("Error during initial tick: %v\n", err)
	}

	tickCount := 0
	configCheckCount := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Scheduler stopped")
			return ctx.Err()
		case <-ticker.C:
			tickCount++
			if err := s.onTick(); err != nil {
				fmt.Printf("Error during tick #%d: %v\n", tickCount, err)
			}
		case <-configTicker.C:
			configCheckCount++
			if s.onConfigReload != nil {
				if err := s.onConfigReload(); err != nil {
					fmt.Printf("Error reloading config (check #%d): %v\n", configCheckCount, err)
				}
			}
		}
	}
}
