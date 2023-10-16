package main

import "fmt"

func printResults(filename string, results []ResultLine) {
	for _, r := range results {
		fmt.Println(r.resultType, r.description)
		fmt.Printf("\t%s:%d:%d to %d:%d\n", filename, r.fl, r.fc, r.ll, r.lc)
	}
}
