package main

import (
	"fmt"
)

type printOptions struct {
	descriptionsIgnored []string
}

func printResults(po printOptions, filename string, results []ResultLine) {
	for _, r := range results {
		if !contains(po.descriptionsIgnored, r.description) {
			fmt.Println(r.resultType, r.description)
			fmt.Printf("\t%s:%d:%d to %d:%d\n", filename, r.fl, r.fc, r.ll, r.lc)
		}
	}
}

func contains(s []string, a string) bool {
	for _, ss := range s {
		if ss == a {
			return true
		}
	}
	return false
}
