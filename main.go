package main

import (
	"net/http"

	"memes-swe-challenge/clients"
	"memes-swe-challenge/log"
	"memes-swe-challenge/scrapper"
)

func main() {

	logger := log.NewLogger()
	pageClient := clients.PageClient{
		Logger: logger,
		Getter: &http.Client{},
	}

	scrapper := scrapper.NewCollector(logger, &pageClient)
	scrapper.OnHTML()
	scrapper.OnRequest()
	scrapper.OnResponse()
	scrapper.OnScraped()
	scrapper.OnVisit()
	scrapper.OnError()
}
