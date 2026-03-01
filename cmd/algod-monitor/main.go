package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/CorvidLabs/go-algod-monitor/internal/config"
	"github.com/CorvidLabs/go-algod-monitor/pkg/health"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	outputFmt string
	interval  int
	oneShot   bool
	maxLag    int
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "algod-monitor",
		Short: "Monitor Algorand node health",
		Long:  "A CLI tool that monitors one or more Algorand algod nodes, reporting health status, round lag, and latency.",
		RunE:  run,
	}

	cmd.Flags().StringVarP(&cfgFile, "config", "c", "", "path to config file (JSON)")
	cmd.Flags().StringVarP(&outputFmt, "output", "o", "text", "output format: text or json")
	cmd.Flags().IntVarP(&interval, "interval", "i", 30, "check interval in seconds")
	cmd.Flags().BoolVar(&oneShot, "once", false, "run a single check and exit")
	cmd.Flags().IntVar(&maxLag, "max-lag", 30, "max acceptable round lag in seconds")

	cmd.AddCommand(versionCmd())
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("algod-monitor v0.1.0")
		},
	}
}

func run(cmd *cobra.Command, args []string) error {
	var cfg config.Config

	if cfgFile != "" {
		var err error
		cfg, err = config.Load(cfgFile)
		if err != nil {
			return err
		}
	} else {
		cfg = config.DefaultConfig()
	}

	// CLI flags override config file values
	if cmd.Flags().Changed("interval") {
		cfg.IntervalSec = interval
	}
	if cmd.Flags().Changed("output") {
		cfg.Output = outputFmt
	}
	if cmd.Flags().Changed("max-lag") {
		cfg.MaxLagSec = maxLag
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	checker := health.NewChecker(
		&http.Client{Timeout: 10 * time.Second},
		time.Duration(cfg.MaxLagSec)*time.Second,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if oneShot {
		results := checker.CheckMultiple(ctx, cfg.Nodes)
		printResults(results, cfg.Output)
		return exitCodeFromResults(results)
	}

	ticker := time.NewTicker(time.Duration(cfg.IntervalSec) * time.Second)
	defer ticker.Stop()

	// Run immediately, then on interval
	results := checker.CheckMultiple(ctx, cfg.Nodes)
	printResults(results, cfg.Output)

	for {
		select {
		case <-ctx.Done():
			fmt.Fprintln(os.Stderr, "\nShutting down.")
			return nil
		case <-ticker.C:
			results := checker.CheckMultiple(ctx, cfg.Nodes)
			printResults(results, cfg.Output)
		}
	}
}

func printResults(results []health.NodeStatus, format string) {
	switch format {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(results)
	default:
		fmt.Printf("\n=== Algod Health Check — %s ===\n", time.Now().Format(time.RFC3339))
		for _, r := range results {
			icon := statusIcon(r.Status)
			fmt.Printf("%s %-30s  status=%-8s  round=%-10d  latency=%s",
				icon, r.Address, r.StatusText, r.LastRound, r.Latency.Round(time.Millisecond))
			if r.Error != "" {
				fmt.Printf("  error=%s", r.Error)
			}
			fmt.Println()
		}
	}
}

func statusIcon(s health.Status) string {
	switch s {
	case health.StatusHealthy:
		return "[OK]"
	case health.StatusDegraded:
		return "[!!]"
	case health.StatusDown:
		return "[XX]"
	default:
		return "[??]"
	}
}

func exitCodeFromResults(results []health.NodeStatus) error {
	var down, degraded []string
	for _, r := range results {
		switch r.Status {
		case health.StatusDown:
			down = append(down, r.Address)
		case health.StatusDegraded:
			degraded = append(degraded, r.Address)
		}
	}
	if len(down) > 0 {
		return fmt.Errorf("nodes down: %s", strings.Join(down, ", "))
	}
	if len(degraded) > 0 {
		return fmt.Errorf("nodes degraded: %s", strings.Join(degraded, ", "))
	}
	return nil
}
