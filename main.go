package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"memes-swe-challenge/clients"
	"memes-swe-challenge/log"
	"memes-swe-challenge/scraper"
)

const (
	exitFail           = 1
	defaultAmountValue = 10
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run() error {
	logger := log.NewLogger()

	if len(os.Args) < 2 {
		return errors.New("expected 'amount' command")
	}
	amountCmdValue := flag.Int("amount", defaultAmountValue, "indicate the amount of images to download")
	flag.Parse()

	logger.Info().Msgf("Starting with amount value: %v", *amountCmdValue)

	pageClient := clients.PageClient{
		Logger: logger,
		Getter: &http.Client{},
	}

	scraper := scraper.NewCollector(logger, &pageClient, *amountCmdValue)
	scraper.ScrapeSite()

	return nil
}
