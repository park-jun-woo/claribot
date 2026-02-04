package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"parkjunwoo.com/claribot/pkg/claude"
	"parkjunwoo.com/claribot/pkg/telegram"

	"gopkg.in/yaml.v3"
)

const Version = "0.2.1"

// Config represents the full configuration
type Config struct {
	Service struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"service"`
	Telegram struct {
		Token        string  `yaml:"token"`
		AllowedUsers []int64 `yaml:"allowed_users"`
	} `yaml:"telegram"`
	Claude struct {
		Timeout int `yaml:"timeout"` // seconds
		Max     int `yaml:"max"`
	} `yaml:"claude"`
}

var bot *telegram.Bot

func main() {
	cfg := loadConfig()

	// Initialize Claude manager
	claude.Init(claude.Config{
		Timeout: time.Duration(cfg.Claude.Timeout) * time.Second,
		Max:     cfg.Claude.Max,
	})

	// Initialize Telegram bot
	if cfg.Telegram.Token != "" && cfg.Telegram.Token != "BOT_TOKEN" {
		var err error
		bot, err = telegram.New(cfg.Telegram.Token)
		if err != nil {
			log.Fatalf("Failed to create telegram bot: %v", err)
		}

		bot.SetHandler(func(msg telegram.Message) {
			log.Printf("[Telegram] %s: %s", msg.Username, msg.Text)
			bot.Send(msg.ChatID, "Echo: "+msg.Text)
		})

		if err := bot.Start(); err != nil {
			log.Fatalf("Failed to start telegram bot: %v", err)
		}
		log.Printf("Telegram bot connected: @%s", bot.Username())
	} else {
		log.Println("Telegram bot disabled (no token configured)")
	}

	// Setup HTTP handler
	http.HandleFunc("/", handleRequest)

	// Start HTTP server in goroutine
	addr := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)
	go func() {
		log.Printf("claribot v%s starting on %s", Version, addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Shutting down...")
	if bot != nil {
		bot.Stop()
	}
	log.Println("Goodbye!")
}

func loadConfig() Config {
	cfg := Config{}
	// defaults
	cfg.Service.Host = "127.0.0.1"
	cfg.Service.Port = 9847
	cfg.Claude.Timeout = 1200
	cfg.Claude.Max = 3

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: cannot get home directory: %v", err)
		return cfg
	}

	configPath := filepath.Join(homeDir, ".claribot", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("Warning: cannot read config file: %v (using defaults)", err)
		return cfg
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Printf("Warning: cannot parse config file: %v (using defaults)", err)
		return cfg
	}

	log.Printf("Config loaded from %s", configPath)
	return cfg
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	args := r.URL.Query().Get("args")

	if args == "" {
		http.Error(w, "missing args parameter", http.StatusBadRequest)
		return
	}

	log.Printf("Received: %s", args)

	// Echo back for now
	fmt.Fprint(w, args)
}
