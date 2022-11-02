package main

import (
	"net/http"

	"memes-swe-challenge/clients"
	"memes-swe-challenge/log"
	"memes-swe-challenge/scraper"
)

func main() {

	logger := log.NewLogger()
	pageClient := clients.PageClient{
		Logger: logger,
		Getter: &http.Client{},
	}

	scraper := scraper.NewCollector(logger, &pageClient)
	scraper.ScrapeSite()
}
