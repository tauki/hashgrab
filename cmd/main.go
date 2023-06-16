//go:build go1.16
// +build go1.16

// hashgrab is a CLI tool that fetches and hashes contents from a list of URLs.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tauki/hashgrab"
)

// Command line flag to set the number of parallel requests.
var parallel = flag.Int("parallel", 10, "limit the number of parallel requests")

func main() {
	// Parse the command line flags.
	flag.Parse()
	if *parallel <= 0 {
		fmt.Println("Number of parallel requests should be greater than 0")
		os.Exit(1)
	}

	// Get the URLs from the command line arguments.
	urls := flag.Args()

	// Validate the provided URLs.
	if len(urls) == 0 {
		fmt.Println("Please provide at least one URL")
		os.Exit(1)
	}

	// Create a new Worker, set the number of workers and start processing the URLs.
	out := hashgrab.New().MaxWorker(*parallel).Run(urls)

	// Loop over the results channel.
	for o := range out {
		if o.Error != nil {
			// Print an error message if fetching or hashing failed.
			fmt.Printf("couldn't fetch %s: %+v\n", o.Url, o.Error)
			continue
		}
		// Print the URL and its hash.
		fmt.Println(o.Url, o.Hash)
	}
}
