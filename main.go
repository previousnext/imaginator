package main

import (
	"fmt"
	"net/url"

	"gopkg.in/alecthomas/kingpin.v2"
	"golang.org/x/sync/errgroup"

	"github.com/previousnext/imaginator/internal/fileutils"
	"github.com/previousnext/imaginator/internal/scrape"
)

var (
	cliSource = kingpin.Arg("source", "Verbose mode.").String()
	cliTarget = kingpin.Arg("target", "Verbose mode.").String()
)

func main() {
	kingpin.Parse()

	err := run(*cliSource, *cliTarget)
	if err != nil {
		panic(err)
	}
}

func run(source, target string) error {
	s, err := url.Parse(source)
	if err != nil {
		return err
	}

	images, err := scrape.Images(s)
	if err != nil {
		return err
	}

	var g errgroup.Group

	for _, image := range images {
		image := image // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			target := fmt.Sprintf("%s%s", target, image.Path)
			fmt.Println("Downloading:", target)
			return fileutils.Download(image.String(), target)
		})
	}
	// Wait for all HTTP fetches to complete.
	return g.Wait()
}
