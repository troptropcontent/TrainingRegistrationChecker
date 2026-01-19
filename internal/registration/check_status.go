package registration

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

const TRAINING_URL_ENV = "TRAINING_URL_ENV"

func IsOpen() (bool, error) {
	c := colly.NewCollector()

	var isOpen bool
	var fetchError error

	c.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		isOpen = !strings.Contains(e.Text, "Candidature fermée.")
	})

	c.OnError(func(r *colly.Response, err error) {
		fetchError = err
	})

	trainingUrl := os.Getenv(TRAINING_URL_ENV)
	if trainingUrl == "" {
		fetchError = fmt.Errorf("environment variable %s is not set", TRAINING_URL_ENV)
	}

	c.Visit(trainingUrl)

	return isOpen, fetchError
}
