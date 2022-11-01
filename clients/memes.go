package clients

import (
	"io/ioutil"
	"net/http"

	"memes-swe-challenge/log"
)

const methodType = "GET"

type PageClient struct {
	Logger *log.Logger
	Getter interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

func (c PageClient) GetImageFromUrl(url string) ([]byte, error) {
	c.Logger.Log().Msgf("calling ur: %v", url)

	req, err := http.NewRequest(methodType, url, nil)

	resp, err := c.Getter.Do(req)
	if err != nil {
		c.Logger.Error().Err(err)
		return []byte{}, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
