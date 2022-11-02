package main

import (
	"flag"
	"net/http"
	"os"

	"memes-swe-challenge/clients"
	"memes-swe-challenge/log"
	"memes-swe-challenge/scrapper"
)

const defaultAmountValue = 10

func main() {
	logger := log.NewLogger()

	if len(os.Args) < 2 {
		logger.Error().Msgf("expected 'amount' command, using default value %v", defaultAmountValue)
	}
	amountCmdValue := flag.Int("amount", defaultAmountValue, "indicate the amount of images to download")
	flag.Parse()

	logger.Info().Msgf("Starting with amount value: %v", *amountCmdValue)

	pageClient := clients.PageClient{
		Logger: logger,
		Getter: &http.Client{},
	}

	scrapper := scrapper.NewCollector(logger, &pageClient, *amountCmdValue)
	scrapper.OnHTML()
	scrapper.OnRequest()
	scrapper.OnResponse()
	scrapper.OnScraped()
	scrapper.OnVisit()
	scrapper.OnError()
}
