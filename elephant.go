package elephant

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	FirstIssue = 41
	issuePage  = "https://golangweekly.com/issues/%d"
	jar, _     = cookiejar.New(nil)
	Client     = newClient()
)

func newClient() *http.Client {
	return &http.Client{Jar: jar}
}

type Issue struct {
	Nbr   int
	Links map[string]*url.URL
	Date  time.Time
}

type Item struct {
	Title        string
	Link         string
	Desc         string
	Tags         []string
	InternalLink string
}

func (i *Issue) Fetch() error {
	doc, err := goquery.NewDocument(fmt.Sprintf(issuePage, i.Nbr))
	if err != nil {
		return err
	}

	doc.Find("table.item").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a")
		if link == nil {
			return
		}
		url, _ := link.Attr("title")
		desc := strings.TrimSpace(s.Find("td.body > div").Text())
		href, _ := link.Attr("href")
		item := &Item{
			InternalLink: href,
			Title:        link.Text(),
			Link:         url,
			Desc:         desc,
			Tags:         []string{},
		}
		spanTag := link.Parent().Find("span.tag")
		if spanTag != nil {
			if classes, ok := spanTag.Attr("class"); ok {
				list := strings.Split(classes, " ")
				for _, c := range list {
					if c != "tag" {
						item.Tags = append(item.Tags, strings.Replace(c, "tag-", "", 1))
					}
				}
			}
		}
		fmt.Printf("%+v\n", item)
	})
	return nil
}
