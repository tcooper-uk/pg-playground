package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tcooper/pg-playground/app/api"
	"github.com/tcooper/pg-playground/app/config"
	"github.com/tcooper/pg-playground/app/db"
	"github.com/tcooper/pg-playground/app/simulator"
)

func main() {
	cfgPath := "config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conns, err := db.Connect(ctx, cfg.Database.PrimaryDSN, cfg.Database.ReplicaDSN)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer conns.Close()
	log.Println("connected to primary and replica")

	metrics := simulator.NewMetrics()
	weights := simulator.NewWeightStore(cfg.Simulator.Weights)
	rate := simulator.NewRateStore(cfg.Simulator.RatePerSecond)

	if cfg.Simulator.Enabled {
		log.Printf("starting simulator: %d workers, %d ops/sec", cfg.Simulator.Workers, cfg.Simulator.RatePerSecond)
		go simulator.Run(ctx, cfg.Simulator, conns.Primary, conns.Replica, metrics, weights, rate)
	}

	addr := fmt.Sprintf(":%d", cfg.API.Port)
	srv := &http.Server{Addr: addr, Handler: api.NewRouter(conns, metrics, weights, rate)}

	go func() {
		log.Printf("API listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down")
	srv.Shutdown(context.Background())
}
