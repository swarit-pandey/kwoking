package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/swarit-pandey/kwoking/internal/ced"
	"github.com/swarit-pandey/kwoking/internal/config"
)

func main() {
	system := flag.String("system", "", "the system that you need to stress test")
	configFilePath := flag.String("config-file", "", "config file for publisher")
	flag.Parse()
	run(*system, *configFilePath)
}

func run(system string, configFilePath string) {
	if system == "" {
		slog.Error("you need to provide system name")
		os.Exit(1)
	}

	if system == "ced" {
		if configFilePath == "" {
			configFilePath = "./ced-internal-dev-config.yaml"
		}
		cedConf := config.NewCEDConfig(configFilePath)

		if err := cedConf.LoadConfig(); err != nil {
			slog.Error("failed to process config", err)
			os.Exit(1)
		}

		if err := cedConf.Unmarshal(); err != nil {
			slog.Error("failed to process config", err)
			os.Exit(1)
		}
		slog.Info("Successfully unmarshalled config", slog.Any("config", cedConf))
		if cedConf == nil {
			slog.Error("Configuration is not initialized properly")
			os.Exit(1)
		}

		cedCore := ced.NewCEDCore(setupCtx(), &cedConf.CED)
		cedCore.Simulate()
	} else {
		slog.Info("only ced is supported as a system as of now")
		os.Exit(1)
	}
}

func setupCtx() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		osSigsChan := make(chan os.Signal, 1)
		signal.Notify(osSigsChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT)

		select {
		case sig := <-osSigsChan:
			slog.Info("Received signal, winding up...", "signal", sig)
			cancel()
		case <-ctx.Done():
			// Context was cancelled elsewhere
		}

		signal.Stop(osSigsChan)
		close(osSigsChan)
	}()

	return ctx
}
