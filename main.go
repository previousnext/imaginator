package main

import (
	"fmt"
	"net/url"

	kingpin "github.com/alecthomas/kingpin/v2"
	"golang.org/x/sync/errgroup"

	"github.com/previousnext/imaginator/internal/fileutils"
	"github.com/previousnext/imaginator/internal/scrape"
)

var (
	cliUser   = kingpin.Flag("user", "Verbose mode.").String()
	cliPass   = kingpin.Flag("pass", "Verbose mode.").String()
	cliSource = kingpin.Arg("source", "Verbose mode.").Required().String()
	cliTarget = kingpin.Arg("target", "Verbose mode.").Required().String()
)

func main() {
	kingpin.Parse()

	err := run(*cliSource, *cliUser, *cliPass, *cliTarget)
	if err != nil {
		panic(err)
	}
}

func run(source, user, pass, target string) error {
	s, err := url.Parse(source)
	if err != nil {
		return err
	}

	images, err := scrape.Images(s, user, pass)
	if err != nil {
		return err
	}

	var g errgroup.Group

	for _, image := range images {
		image := image // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			target := fmt.Sprintf("%s%s", target, image.Path)
			fmt.Println("Downloading:", target)

			err := fileutils.Download(image.String(), user, pass, target)
			if err != nil {
				fmt.Printf("Failed to download: %w\n", err)
			}

			return nil
		})
	}
	// Wait for all HTTP fetches to complete.
	return g.Wait()
}
