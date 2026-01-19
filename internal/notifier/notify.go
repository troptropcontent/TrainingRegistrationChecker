package notifier

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const BASE_NTFY_URL = "https://ntfy.sh/"

type Channel string

const (
	Success Channel = "NTFY_SUCCESS_CHANNEL"
	Logs    Channel = "NTFY_LOGS_CHANNEL"
)

func Notify(message string, channel Channel) error {
	// Always log to console for debugging
	log.Printf("[%s] %s", channel, message)

	ntfyChannel := os.Getenv(string(channel))
	if ntfyChannel == "" {
		return fmt.Errorf("environment variable %s is not set", channel)
	}

	ntfyURL := BASE_NTFY_URL + ntfyChannel
	_, err := http.Post(ntfyURL, "text/plain", strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("unable to send notification: %v", err)
	}

	return nil
}
