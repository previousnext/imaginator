package scrape

import (
	"encoding/base64"
	"log"
	"fmt"
	"net/url"
	"strings"

	"github.com/gocolly/colly"
)

// Images extracted from a path.
func Images(source *url.URL, user, pass string) ([]*url.URL, error) {
	var images []*url.URL

	c := colly.NewCollector(
		colly.AllowedDomains(source.Host),
	)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Authorization", "Basic "+basicAuth(user, pass))
	})

	c.OnHTML("source[srcset]", func(e *colly.HTMLElement) {
		srcset := e.Attr("srcset")

		links, err := splitSrcSet(source, srcset)
		if err != nil {
			log.Println(err)
			return
		}

		images = append(images, links...)
	})

	c.OnHTML("img[src]", func(e *colly.HTMLElement) {
		image, err := url.Parse(fmt.Sprintf("%s://%s%s", source.Scheme, source.Host, e.Attr("src")))
		if err != nil {
			log.Println(err)
			return
		}

		images = append(images, image)
	})

	err := c.Visit(source.String())
	if err != nil {
		return images, err
	}

	var originals []*url.URL

	for _, image := range images {
		original, err := getOriginalFromStyle(image)
		if err != nil {
			log.Println(err)
		}

		originals = append(originals, original)
	}

	images = append(images, originals...)

	return images, nil
}

func splitSrcSet(source *url.URL, set string) ([]*url.URL, error) {
	var images []*url.URL

	for _, image := range strings.Split(set, ", ") {
		sl := strings.Split(image, " ")

		if len(sl) == 2 {
			image, err := url.Parse(fmt.Sprintf("%s://%s%s", source.Scheme, source.Host, sl[0]))
			if err != nil {
				return images, err
			}

			images = append(images, image)
		}
	}

	return images, nil
}

func getOriginalFromStyle(link *url.URL) (*url.URL, error) {
	split := strings.Split(link.Path, "/")
	newPath := append(split[:4], split[7:]...)

	original := &url.URL{
		Scheme: link.Scheme,
		Host: link.Host,
		Path: strings.Join(newPath, "/"),
	}

	return  original, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
