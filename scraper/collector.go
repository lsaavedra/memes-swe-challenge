package scraper

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"

	"memes-swe-challenge/log"
)

const (
	imagesPath  = "images"
	imagesLimit = 10
	pageUri     = "https://icanhas.cheezburger.com/"
)

type (
	Scraper struct {
		logger         *log.Logger
		collyCollector *colly.Collector
		pageClient     pageClient
		memes          []Meme
	}
	pageClient interface {
		GetImageFromUrl(url string) ([]byte, error)
	}
)

func NewCollector(logger *log.Logger, pageClient pageClient) *Scraper {
	collector := colly.NewCollector()
	collector.SetRequestTimeout(120 * time.Second)
	return &Scraper{
		logger:         logger,
		collyCollector: collector,
		pageClient:     pageClient,
		memes:          make([]Meme, 0),
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

	s.collyCollector.OnRequest(func(r *colly.Request) {
		s.logger.Info().Msgf("Visiting: %v", r.URL)
	})

	s.collyCollector.OnResponse(func(r *colly.Response) {
		s.logger.Info().Msgf("Got a response from %v", r.Request.URL)
	})

	s.collyCollector.OnError(func(r *colly.Response, e error) {
		s.logger.Error().Err(e)
	})

	s.collyCollector.OnScraped(func(r *colly.Response) {
		s.logger.Info().Msgf("Finished %v", r.Request.URL)
		if err := os.Mkdir(imagesPath, os.ModePerm); err != nil {
			s.logger.Error().Err(err)
		}
		for idx, img := range s.memes {
			if idx == imagesLimit {
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

	})

	s.collyCollector.Visit(pageUri)
}

func buildImageName(value int) string {
	return fmt.Sprintf("%v/%v.jpg", imagesPath, strconv.Itoa(value+1))
}
