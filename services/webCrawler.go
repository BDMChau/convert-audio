package services

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/llimllib/loglevel"
	"golang.org/x/net/html"
)

var visitedUrl = []string{}
var MaxDepth = 2

type Link struct {
	url   string
	text  string
	depth int
}

type HttpError struct {
	original string
}

func LinkReader(resp *http.Response, depth int) []Link {
	page := html.NewTokenizer(resp.Body)
	links := []Link{}

	var start *html.Token
	var text string

	for {
		_ = page.Next()
		token := page.Token()

		if token.Type == html.ErrorToken {
			break
		}

		if start == nil && token.Type == html.TextToken {
			fmt.Println("RespondData:", token.Data)
			text = fmt.Sprintf("%s %s", text, token.Data)
		}

		switch token.Type {
		case html.StartTagToken:
			if len(token.Attr) > 0 {
				start = &token
			}
		case html.EndTagToken:
			if start == nil {
				fmt.Println("Link end found without start: %s", text)
				continue
			}
			link := NewLink(*start, text, depth)
			if link.Valid() {
				links = append(links, link)
				log.Debugf("Link Found %v", link)
			}

			start = nil
			text = ""
		}

	}

	log.Debug(links)
	return links
}

func NewLink(tag html.Token, text string, depth int) Link {
	link := Link{text: strings.TrimSpace(text), depth: depth}

	for i := range tag.Attr {
		if tag.Attr[i].Key == "href" {
			link.url = strings.TrimSpace(tag.Attr[i].Val)
		}
	}

	return link
}

func (self Link) String() string {
	spacer := strings.Repeat("/t", self.depth)
	return fmt.Sprintf("%s %s (%d) - %s", spacer, self.text, self.depth, self.url)
}

func (self Link) Valid() bool {
	if self.depth >= MaxDepth {
		return false
	}
	if len(self.text) == 0 {
		return false
	}
	if len(self.url) == 0 || strings.Contains(strings.ToLower(self.url), "javascript") {
		return false
	}

	return true
}

func (self HttpError) Error() string {
	return self.original
}

func downloader(url string) (resp *http.Response, err error) {
	log.Debugf("Downloading %s", url)

	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > 299 {
		err = HttpError{fmt.Sprintf("Error: ", resp, url)}
		return nil, err
	}

	return resp, nil
}

func recursiveDownloader(url string, depth int) {
	resp, err := downloader(url)
	if err != nil {
		log.Error(err)
		return
	}

	links := LinkReader(resp, depth)
	defer resp.Body.Close()

	for _, link := range links {
		fmt.Println(link)

		plusDepth := depth + 1
		if plusDepth < MaxDepth {
			recursiveDownloader(link.url, plusDepth)
		}
	}
}

func Crawler() {
	log.SetPriorityString("info")
	log.SetPrefix("Crawler:::")
	log.Debug(os.Args)

	if len(os.Args) < 2 {
		log.Fatalln("Missing Url arg")
	}

	recursiveDownloader(os.Args[1], 0)
}
