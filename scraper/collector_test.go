package scraper

import (
	"errors"
	"testing"

	"memes-swe-challenge/log"
)

type mockPageClient struct {
	fail bool
}

func (m *mockPageClient) GetImageFromUrl(url string) ([]byte, error) {
	if m.fail {
		return []byte{}, errors.New("fail getting image")
	}
	return []byte("test"), nil
}

func TestScraper_getImage(t *testing.T) {
	type args struct {
		url string
		idx int
	}
	tests := []struct {
		name       string
		logger     *log.Logger
		pageClient mockPageClient
		args       args
		wantErr    bool
	}{
		{
			name:       "getImage - Fail by client",
			logger:     log.NewLogger(),
			pageClient: mockPageClient{fail: true},
			args:       args{url: "https://i.chzbgr.com/thumb800/18237957/hB702534D/of-a-cat-hand-walked-in-on-my-husband-using-our-cat-as-a-mobile-check-deposit-background-ming-ho", idx: 1},
			wantErr:    true,
		},
		{
			name:       "getImage - Fail by directory error",
			logger:     log.NewLogger(),
			pageClient: mockPageClient{fail: false},
			args:       args{url: "https://i.chzbgr.com/thumb800/18237957/hB702534D/of-a-cat-hand-walked-in-on-my-husband-using-our-cat-as-a-mobile-check-deposit-background-ming-ho", idx: 1},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := NewCollector(tt.logger, &tt.pageClient, 1, 1)
			if err := collector.getImage(tt.args.url, tt.args.idx); (err != nil) != tt.wantErr {
				t.Errorf("getImage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
