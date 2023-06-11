package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tauki/hashgrab"
)

var parallel = flag.Int("parallel", 10, "limit the number of parallel requests")

func main() {
	flag.Parse()
	if *parallel <= 0 {
		fmt.Println("Number of parallel requests should be greater than 0")
		os.Exit(1)
	}
	urls := flag.Args()
	if len(urls) == 0 {
		fmt.Println("Please provide at least one URL")
		os.Exit(1)
	}
	out := hashgrab.New().MaxWorker(*parallel).Run(urls)
	for o := range out {
		if o.Error != nil {
			fmt.Printf("couldn't fetch %s: %+v\n", o.Url, o.Error)
			continue
		}
		fmt.Println(o.Url, o.Hash)
	}
}
