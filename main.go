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
	exitFail            = 1
	defaultAmountValue  = 10
	defaultThreadsValue = 1
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

func run() error {
	logger := log.NewLogger()

	if len(os.Args) < 3 {
		return errors.New("expected 'amount' or 'threads' command")
	}

	amountCmdValue := flag.Int("amount", defaultAmountValue, "indicate the amount of images to download")
	threadsCmdValue := flag.Int("threads", defaultThreadsValue, "indicate how many threads will run")
	flag.Parse()

	if *threadsCmdValue < 1 || *threadsCmdValue > 5 {
		logger.Fatal().Msgf("threads command line value needs to be in range from 1 to 5")
		return errors.New("threads command line value needs to be in range from 1 to 5")
	}

	logger.Info().Msgf("Starting with amount value: %v and threads: %v", *amountCmdValue, *threadsCmdValue)

	pageClient := clients.PageClient{
		Logger: logger,
		Getter: &http.Client{},
	}

	scraper := scraper.NewCollector(logger, &pageClient, *amountCmdValue, *threadsCmdValue)
	scraper.ScrapeSite()

	return nil
}
