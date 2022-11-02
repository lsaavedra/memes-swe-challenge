package scraper

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"

	"memes-swe-challenge/clients"
	"memes-swe-challenge/log"
)

const (
	imagesPath = "images"
	pageUri    = "https://icanhas.cheezburger.com/"
)

type Scraper struct {
	logger         *log.Logger
	collyCollector *colly.Collector
	pageClient     *clients.PageClient
	memes          []Meme
	imagesLimit    int
}

func NewCollector(logger *log.Logger, pageClient *clients.PageClient, imagesLimit int) *Scraper {
	collector := colly.NewCollector()
	collector.SetRequestTimeout(120 * time.Second)
	return &Scraper{
		logger:         logger,
		collyCollector: collector,
		pageClient:     pageClient,
		memes:          make([]Meme, 0),
		imagesLimit:    imagesLimit,
	}
}

func (s *Scraper) ScrapeSite() {
	s.collyCollector.OnHTML("div.mu-post.mu-thumbnail.resp-media-wrap", func(e *colly.HTMLElement) {
		item := Meme{}
		item.Title = e.ChildAttr("img", "title")
		item.Src = e.ChildAttr("img", "src")
		item.DataSrc = e.ChildAttr("img", "data-src")
		s.memes = append(s.memes, item)
	})

	s.collyCollector.OnHTML("[aria-label=\"Go to next page\"]", func(e *colly.HTMLElement) {
		if len(s.memes) < s.imagesLimit {
			nextPage := e.Request.AbsoluteURL(e.Attr("href"))
			s.logger.Info().Msgf("Going to next page %v", nextPage)
			s.collyCollector.Visit(nextPage)
		}
	})

	s.collyCollector.OnRequest(func(r *colly.Request) {
		s.logger.Info().Msgf("Visiting: %v", r.URL)
	})

	s.collyCollector.OnResponse(func(r *colly.Response) {
		s.logger.Info().Msgf("Got a response from %v", r.Request.URL)
	})

	s.collyCollector.OnError(func(r *colly.Response, e error) {
		s.logger.Error().Err(e)
	})

	s.collyCollector.Visit(pageUri)

	s.collyCollector.Wait()

	if err := os.Mkdir(imagesPath, os.ModePerm); err != nil {
		s.logger.Error().Err(err)
	}
	s.logger.Info().Msgf("Total memes %v", len(s.memes))
	for idx, img := range s.memes {
		if idx == s.imagesLimit {
			break
		}
		bytesFile, err := s.pageClient.GetImageFromUrl(img.DataSrc)
		if err != nil {
			s.logger.Error().Err(err)
		}
		if err := os.WriteFile(buildImageName(idx), bytesFile, 0664); err != nil {
			s.logger.Error().Err(err)
		}
	}
	s.logger.Info().Msg("Finished saving images in directory")
}

func buildImageName(value int) string {
	return fmt.Sprintf("%v/%v.jpg", imagesPath, strconv.Itoa(value+1))
}
