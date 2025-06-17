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
}

func Usage() {
	fmt.Fprintf(os.Stderr, "\n🐳 Docker Notifier - Monitor container logs and send alerts to Telegram\n\n")
	fmt.Fprintf(os.Stderr, "Usage: dockernotify [flags]\n\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	pflag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExample:\n  dockernotify --interval 30 --telegram-token \"xxx\" --telegram-chat-id \"-xxx\" --error-pattern \"ERROR\" --error-pattern \"FATAL\"\n\n")
}

func Parse() (*Config, error) {
	interval := pflag.Int("interval", 5, "Log polling interval in seconds")
	labelEnable := pflag.Bool("label-enable", false, "Enable label filter: com.andvarfolomeev.dockernotifier.enable=true")
	telegramToken := pflag.String("telegram-token", "", "Telegram Bot API token")
	telegramChatID := pflag.String("telegram-chat-id", "", "Target chat ID")
	debug := pflag.Bool("debug", false, "Enable debug logging")

	var errorPatterns []string
	pflag.StringSliceVar(&errorPatterns, "error-pattern", []string{"ERROR"}, "Regex pattern for matching error lines (can be used multiple times)")

	help := pflag.BoolP("help", "h", false, "Display help information")

	pflag.Usage = Usage

	pflag.Parse()

	if *help {
		pflag.Usage()
		return nil, ErrHelpRequested
	}

	if *telegramToken == "" {
		return nil, ErrMissingArg("--telegram-token")
	}

	if *telegramChatID == "" {
		return nil, ErrMissingArg("--telegram-chat")
	}

	config := &Config{
		Interval:       *interval,
		LabelEnable:    *labelEnable,
		TelegramToken:  *telegramToken,
		TelegramChatID: *telegramChatID,
		ErrorPatterns:  errorPatterns,
		Debug:          *debug,
	}

	return config, nil
}
