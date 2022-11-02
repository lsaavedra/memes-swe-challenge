package main

import (
	"flag"
	"net/http"
	"os"

	"memes-swe-challenge/clients"
	"memes-swe-challenge/log"
	"memes-swe-challenge/scraper"
)

const (
	defaultAmountValue  = 10
	defaultThreadsValue = 1
)

func main() {
	logger := log.NewLogger()

	if len(os.Args) < 3 {
		logger.Error().Msg("expected 'amount' or 'threads' command, using default values")
	}
	amountCmdValue := flag.Int("amount", defaultAmountValue, "indicate the amount of images to download")
	threadsCmdValue := flag.Int("threads", defaultThreadsValue, "indicate how many threads will run")
	flag.Parse()

	if *threadsCmdValue < 1 || *threadsCmdValue > 5 {
		logger.Fatal().Msgf("threads command line value needs to be in range from 1 to 5")
	}

	logger.Info().Msgf("Starting with amount value: %v and threads: %v", *amountCmdValue, *threadsCmdValue)

	pageClient := clients.PageClient{
		Logger: logger,
		Getter: &http.Client{},
	}

	scraper := scraper.NewCollector(logger, &pageClient, *amountCmdValue, *threadsCmdValue)
	scraper.ScrapeSite()
}
