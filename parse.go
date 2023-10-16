package main

import (
	"io"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type ResultLine struct {
	resultType  string
	description string
	fl          int
	fc          int
	ll          int
	lc          int
}

func parse(r io.Reader) ([]ResultLine, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	results := []ResultLine{}
	for _, label := range []string{".error", ".warning", ".info"} {
		doc.Find(label).Each(func(_ int, s *goquery.Selection) {
			description := s.Find("p span").First().Text()
			resultType := s.Find("p strong").Text()
			fl, err := strconv.Atoi(s.Find("p.location span.first-line").Text())
			if err != nil {
				return
			}
			fc, err := strconv.Atoi(s.Find("p.location span.first-col").Text())
			if err != nil {
				return
			}
			ll, err := strconv.Atoi(s.Find("p.location span.last-line").Text())
			if err != nil {
				return
			}
			lc, err := strconv.Atoi(s.Find("p.location span.last-col").Text())
			if err != nil {
				return
			}
			results = append(results, ResultLine{
				resultType:  resultType,
				description: description,
				fl:          fl,
				fc:          fc,
				ll:          ll,
				lc:          lc,
			})
		})
	}
	return results, nil
}
