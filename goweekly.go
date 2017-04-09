package elephant

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	// FirstGoWeeklyIssue is the first available golang weekly issue #
	FirstGoWeeklyIssue = 41
	GoWeeklyIssuePage  = "https://golangweekly.com/issues/%d"
)

// Issue is the representation of a Go Weekly issue
type Issue struct {
	// Nbr is the issue number
	Nbr int
	// Items is the number of items found in the issue
	Items []*Item
	// Date is the issue date
	Date time.Time
	// Doc represents the HTML document of the issue
	Doc *goquery.Document
}

// Fetch retrieves the issue from the internet.
func (it *Issue) Fetch() (err error) {
	it.Doc, err = goquery.NewDocument(fmt.Sprintf(GoWeeklyIssuePage, it.Nbr))
	return err
}

// Parse parses the issue document and extracts the items.
func (it *Issue) Parse() error {
	if it == nil {
		return fmt.Errorf("failed to parse the issue (nil)")
	}
	if metadata := it.Doc.Find("div.issuemetadata").Text(); metadata != "" {
		// Issue 99 — March  3, 2016
		chunks := strings.Split(metadata, "—")
		if len(chunks) > 1 {
			dateStr := strings.TrimSpace(chunks[1])
			var err error
			it.Date, err = time.Parse("January 2, 2006", dateStr)
			if err != nil {
				fmt.Println("failed parsing date", dateStr, err)
			}
		}
	}
	it.Items = []*Item{}
	it.Doc.Find("table.item").Each(func(i int, s *goquery.Selection) {
		// classes, _ := s.Attr("class")
		// fmt.Printf("-- %s\n", classes)
		// fmt.Println(s.HasClass("section-jobs"))
		if s.HasClass("section-jobs") {
			return
		}
		link := s.Find("a")
		if link == nil {
			return
		}
		url, _ := link.Attr("title")
		desc := strings.TrimSpace(s.Find("td.body > div").Text())
		sourceEl := s.Find("td.source > div")
		if sourceEl.Has("span.tag-sponsored").Length() > 0 {
			return
		}
		source := strings.TrimSpace(sourceEl.Text())

		href, _ := link.Attr("href")
		item := &Item{
			InternalLink: href,
			Title:        link.Text(),
			Source:       source,
			Link:         url,
			Desc:         desc,
			Tags:         []string{},
			Date:         it.Date,
		}
		spanTag := link.Parent().Find("span.tag")
		if spanTag != nil {
			if classes, ok := spanTag.Attr("class"); ok {
				list := strings.Split(classes, " ")
				for _, c := range list {
					if c != "tag" {
						tag := strings.Replace(c, "tag-", "", 1)
						item.Tags = append(item.Tags, tag)
					}
				}
			}
		}
		it.Items = append(it.Items, item)
	})
	return nil
}
