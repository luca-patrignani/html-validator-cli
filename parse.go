package main

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)


func parse(r io.Reader) error {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return err
	}
	for _, label := range []string{".error", ".warning", ".info"} {
		doc.Find(label).Each(func(_ int, s *goquery.Selection) {
			description := s.Find("p span").First().Text()
			resultType := s.Find("p strong").Text()
			fmt.Println(resultType, ":", description)
			fl := s.Find("p.location span.first-line").Text()
			fc := s.Find("p.location span.first-col").Text()
			ll := s.Find("p.location span.last-line").Text()
			lc := s.Find("p.location span.last-col").Text()
			fmt.Printf("\tFrom %s:%s to %s:%s\n", fl, fc, ll, lc)
		})
	}
	return nil
}