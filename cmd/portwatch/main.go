package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/schedule"
	"github.com/user/portwatch/internal/store"
)

func main() {
	cfgPath := flag.String("config", "portwatch.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ports, err := scanner.ParsePortRange(cfg.Ports)
	if err != nil {
		log.Fatalf("ports: %v", err)
	}

	st, err := store.New(cfg.StorePath)
	if err != nil {
		log.Fatalf("store: %v", err)
	}

	sc := scanner.New(cfg.Timeout)
	notifier := alert.NewMultiNotifier(alert.NewConsoleNotifier())

	runner := schedule.New(schedule.Config{
		Hosts:    cfg.Hosts,
		Ports:    ports,
		Interval: cfg.Interval,
	}, sc, st, notifier)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("portwatch started")
	if err := runner.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("runner: %v", err)
	}
	log.Println("portwatch stopped")
}
