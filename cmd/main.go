package main

import (
	"flag"
	cfg "github.com/17neverends/buildcast/internal/config"
	"github.com/17neverends/buildcast/internal/core"
	"log"
	"os"
	"path/filepath"
)

func main() {
	configPath := flag.String("config", "config.json", "Path to config file")
	serviceFlag := flag.String("service", "", "Service name to append to deploy path")
	flag.Parse()

	config, err := cfg.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	originalEnv, err := os.ReadFile(config.FrontendEnvPath)
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	log.Println("Checking servers availability...")
	availableServers := core.CheckServers(config.Servers)
	if len(availableServers) == 0 {
		log.Fatal("No servers available, exiting")
	}

	for _, server := range availableServers {
		log.Printf("Processing server: %s", server.IP)

		modifiedEnv := core.ModifyEnv(originalEnv, server.Host, config.EnvHost)
		if err := os.WriteFile(config.FrontendEnvPath, modifiedEnv, 0644); err != nil {
			log.Printf("Error writing modified .env: %v", err)
			continue
		}

		log.Println("Running build command...")
		if err := core.RunCommand(config.MainCmd); err != nil {
			log.Printf("Build failed: %v", err)
			continue
		}

		deployPath := filepath.Join(server.Path, *serviceFlag)
		log.Printf("Deploying to %s...", deployPath)
		if err := core.DeployToServer(server, deployPath, config.BuildOutput); err != nil {
			log.Printf("Deploy failed: %v", err)
		}

		if err := os.WriteFile(config.FrontendEnvPath, originalEnv, 0644); err != nil {
			log.Printf("Error restoring .env: %v", err)
		}
	}

	log.Println("Success!")
}
