package scraper

import (
	"fmt"
	"os"
	"strconv"
	"sync"
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
	threads        int
}

func NewCollector(logger *log.Logger, pageClient *clients.PageClient, imagesLimit int, threadsValue int) *Scraper {
	collector := colly.NewCollector()
	collector.SetRequestTimeout(120 * time.Second)
	return &Scraper{
		logger:         logger,
		collyCollector: collector,
		pageClient:     pageClient,
		memes:          make([]Meme, 0),
		imagesLimit:    imagesLimit,
		threads:        threadsValue,
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

	s.getImgAndSave()
}

func (s *Scraper) getImgAndSave() {
	if err := os.Mkdir(imagesPath, os.ModePerm); err != nil {
		s.logger.Error().Err(err)
	}
	channel := make(chan MemeToStore)
	var wg sync.WaitGroup
	wg.Add(s.threads)

	for i := 0; i < s.threads; i++ {
		go func(c chan MemeToStore) {
			for {
				value, more := <-c
				if more == false {
					wg.Done()
					return
				}
				err := s.getImage(value.imageUrl, value.id)
				if err != nil {
					s.logger.Error().Err(err)
				}
			}
		}(channel)
	}
	for idx, meme := range s.memes {
		channel <- MemeToStore{imageUrl: meme.DataSrc, id: idx}
	}
	close(channel)
	wg.Wait()
}

func (s *Scraper) getImage(url string, idx int) error {
	bytesFile, err := s.pageClient.GetImageFromUrl(url)
	if err != nil {
		s.logger.Error().Err(err)
		return err
	}
	if err := os.WriteFile(buildImageName(idx), bytesFile, 0664); err != nil {
		s.logger.Error().Err(err)
		return err
	}
	return nil
}

func buildImageName(value int) string {
	return fmt.Sprintf("%v/%v.jpg", imagesPath, strconv.Itoa(value+1))
}
