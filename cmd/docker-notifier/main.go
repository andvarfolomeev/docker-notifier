package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Define command line flags
	interval := flag.Int("interval", 5, "Log polling interval in seconds")
	labelEnable := flag.Bool("label-enable", false, "Enable label filter: com.andvarfolomeev.dockernotify.enable=true")
	telegramToken := flag.String("telegram-token", "", "Telegram Bot API token")
	telegramChatID := flag.String("telegram-chat-id", "", "Target chat ID")
	debug := flag.Bool("debug", false, "Enable debug logging")
	cleanup := flag.Bool("cleanup", false, "Clear saved log offsets")

	// Support for multiple error patterns
	var errorPatterns stringSliceFlag
	flag.Var(&errorPatterns, "error-pattern", "Regex pattern for matching error lines (can be used multiple times)")

	flag.Parse()

	// If no error patterns provided, use default
	if len(errorPatterns) == 0 {
		errorPatterns = append(errorPatterns, "ERROR")
	}

	// Validate required parameters
	if *telegramToken == "" {
		fmt.Println("Error: --telegram-token is required")
		flag.Usage()
		os.Exit(1)
	}

	if *telegramChatID == "" {
		fmt.Println("Error: --telegram-chat-id is required")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("🐳 Docker Notifier starting...")
	fmt.Printf("Polling interval: %d seconds\n", *interval)
	fmt.Printf("Label filtering: %v\n", *labelEnable)
	fmt.Printf("Error patterns: %v\n", errorPatterns)
	fmt.Printf("Debug mode: %v\n", *debug)
	fmt.Printf("Cleanup mode: %v\n", *cleanup)

	// TODO: Initialize components
	// - Docker client
	// - Telegram client
	// - Watcher service
	
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Block until we receive a signal
	sig := <-sigChan
	fmt.Printf("Received signal %v, shutting down...\n", sig)
	
	// TODO: Perform cleanup
}

// stringSliceFlag is a custom flag type to handle multiple values for the same flag
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}