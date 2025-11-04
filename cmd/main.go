//go:build go1.16
// +build go1.16

// hashgrab is a CLI tool that fetches and hashes contents from a list of URLs.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/tauki/hashgrab"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	os.Exit(run(ctx, os.Stdout, os.Stderr, os.Args[1:]))
}

func run(ctx context.Context, stdout, stderr io.Writer, args []string) int {
	fs := flag.NewFlagSet("hashgrab", flag.ContinueOnError)
	fs.SetOutput(stderr)
	parallel := fs.Int("parallel", 10, "limit the number of parallel requests")
	hashAlg := fs.String("hash", "sha256", "hash algorithm to use (sha256 or md5)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *parallel <= 0 {
		fmt.Fprintln(stderr, "Number of parallel requests should be greater than 0")
		return 1
	}

	hasher, err := hashgrab.NewHasherFromName(*hashAlg)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	urls := fs.Args()
	if len(urls) == 0 {
		fmt.Fprintln(stderr, "Please provide at least one URL")
		return 1
	}

	worker := hashgrab.New().MaxWorker(*parallel).Hasher(hasher)
	out := worker.RunContext(ctx, urls)

	for o := range out {
		if o.Error != nil {
			fmt.Fprintf(stderr, "couldn't fetch %s: %v\n", o.Url, o.Error)
			continue
		}
		fmt.Fprintf(stdout, "%s %s\n", o.Url, o.Hash)
	}

	return 0
}
