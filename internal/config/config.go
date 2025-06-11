package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type Config struct {
	Interval       int
	LabelEnable    bool
	TelegramToken  string
	TelegramChatID string
	ErrorPatterns  []string
	Debug          bool
	Cleanup        bool
}

func Parse() (*Config, error) {
	interval := pflag.Int("interval", 5, "Log polling interval in seconds")
	labelEnable := pflag.Bool("label-enable", false, "Enable label filter: com.andvarfolomeev.dockernotify.enable=true")
	telegramToken := pflag.String("telegram-token", "", "Telegram Bot API token")
	telegramChatID := pflag.String("telegram-chat-id", "", "Target chat ID")
	debug := pflag.Bool("debug", false, "Enable debug logging")
	cleanup := pflag.Bool("cleanup", false, "Clear saved log offsets")

	var errorPatterns []string
	pflag.StringSliceVar(&errorPatterns, "error-pattern", []string{"ERROR"}, "Regex pattern for matching error lines (can be used multiple times)")

	help := pflag.BoolP("help", "h", false, "Display help information")

	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\n🐳 Docker Notifier - Monitor container logs and send alerts to Telegram\n\n")
		fmt.Fprintf(os.Stderr, "Usage: dockernotify [flags]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		pflag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n  dockernotify --interval 30 --telegram-token \"xxx\" --telegram-chat-id \"-xxx\" --error-pattern \"ERROR\" --error-pattern \"FATAL\"\n\n")
	}

	pflag.Parse()

	if *help {
		pflag.Usage()
		os.Exit(0)
	}

	var missingRequired bool

	if *telegramToken == "" {
		fmt.Println("Error: --telegram-token is required")
		missingRequired = true
	}

	if *telegramChatID == "" {
		fmt.Println("Error: --telegram-chat-id is required")
		missingRequired = true
	}

	if missingRequired {
		fmt.Println()
		pflag.Usage()
		return nil, fmt.Errorf("missing required arguments")
	}

	config := &Config{
		Interval:       *interval,
		LabelEnable:    *labelEnable,
		TelegramToken:  *telegramToken,
		TelegramChatID: *telegramChatID,
		ErrorPatterns:  errorPatterns,
		Debug:          *debug,
		Cleanup:        *cleanup,
	}

	return config, nil
}
