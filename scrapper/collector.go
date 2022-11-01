package scrapper

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
	imagesPath  = "images"
	imagesLimit = 10
)

type Scrapper struct {
	logger         *log.Logger
	collyCollector *colly.Collector
	pageClient     *clients.PageClient
	memes          []Meme
}

func NewCollector(logger *log.Logger, pageClient *clients.PageClient) *Scrapper {
	collector := colly.NewCollector()
	collector.SetRequestTimeout(120 * time.Second)
	return &Scrapper{
		logger:         logger,
		collyCollector: collector,
		pageClient:     pageClient,
		memes:          make([]Meme, 0),
	}
}

func (s *Scrapper) OnHTML() {
	s.collyCollector.OnHTML("div.mu-post.mu-thumbnail.resp-media-wrap", func(e *colly.HTMLElement) {
		item := Meme{}
		item.Title = e.ChildAttr("img", "title")
		item.Src = e.ChildAttr("img", "src")
		item.DataSrc = e.ChildAttr("img", "data-src")
		s.memes = append(s.memes, item)
	})
}
func (s *Scrapper) OnRequest() {
	s.collyCollector.OnRequest(func(r *colly.Request) {
		s.logger.Info().Msgf("Visiting: %v", r.URL)
	})
}
func (s *Scrapper) OnResponse() {
	s.collyCollector.OnResponse(func(r *colly.Response) {
		s.logger.Info().Msgf("Got a response from %v", r.Request.URL)
	})
}
func (s *Scrapper) OnError() {
	s.collyCollector.OnError(func(r *colly.Response, e error) {
		s.logger.Error().Err(e)
	})
}
func (s *Scrapper) OnScraped() {
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
}
func (s *Scrapper) OnVisit() {
	s.collyCollector.Visit("https://icanhas.cheezburger.com/")
}

func buildImageName(value int) string {
	return fmt.Sprintf("%v/%v.jpg", imagesPath, strconv.Itoa(value+1))
}
