package surfio

import (
	"github.com/PuerkitoBio/goquery"
	srctype "github.com/chaewonkong/thumbnail-scraper/src_type"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/browser"
	"gopkg.in/headzoo/surf.v1"
)

type ThumbnailScraper struct {
	TargetURL string `json:"target_url"`
	Bow       *browser.Browser
	Images    []Image
}

// New constructor
func New(url string) *ThumbnailScraper {
	bow := surf.NewBrowser()
	bow.AddRequestHeader("Accept", "text/html")
	bow.AddRequestHeader("Accept-Charset", "utf8")
	bow.SetUserAgent(agent.Chrome())

	return &ThumbnailScraper{
		TargetURL: url,
		Bow:       bow,
	}
}

// GetThumbnailImages get all thumbnail images
func (th *ThumbnailScraper) GetThumbnailImages() (images []Image, err error) {
	err = th.Bow.Open(th.TargetURL)
	if err != nil {
		return
	}

	th.getOpenGraphImageFromMetaProperty()

	if len(th.Images) == 0 {
		th.getAnyImageFromImgTag()
	}

	if len(th.Images) == 0 {
		err = NewImgNotFoundError()
		return
	}

	images = th.Images

	return
}

func (th *ThumbnailScraper) getAnyImageFromImgTag() {
	th.Bow.Dom().Find("img").Each(func(i int, s *goquery.Selection) {
		if result, ok := s.Attr("src"); ok {
			image := Image{
				Url:     result,
				SrcType: srctype.IMG_TAG,
			}
			th.Images = append(th.Images, image)
		}
	})
}

// get og:image content from site meta
func (th *ThumbnailScraper) getOpenGraphImageFromMetaProperty() {
	th.Bow.Dom().Find("meta").Each(func(i int, s *goquery.Selection) {
		if result, available := s.Attr("property"); available && result == srctype.OG {
			imgURL, available := s.Attr("content")
			if available {
				image := Image{
					Url:     imgURL,
					SrcType: srctype.OG,
				}
				th.Images = append(th.Images, image)
			}
		}
	})
}
