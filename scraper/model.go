package scraper

type Meme struct {
	Url     string
	Src     string
	DataSrc string
	Title   string
	Width   int
	Height  int
}

type MemeToStore struct {
	imageUrl string
	id       int
}
