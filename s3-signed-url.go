package main

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"fmt"
	"flag"
	"time"
	"os"
)

func main() {
	bucket := flag.String("bucket", "", "bucket")
	key := flag.String("key", "", "key")
	region := flag.String("region", aws.USEast.Name, "region")
	expiration := flag.String("expiration", "1h", "expiration in the Go duration format.  A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \"300ms\", \"-1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")

	flag.Parse()

	if *bucket == "" || *key == "" {
		fmt.Printf("bucket and key are required\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	auth, err := aws.GetAuth("", "")

	if err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stderr, "auth: %#v\n", auth)

	reg, found := aws.Regions[*region]

	if ! found {
		panic(fmt.Sprintf("invalid region: %s", *region))
	}

	s3c := s3.New(auth, reg)
	buck := s3c.Bucket(*bucket)

	dur, err := time.ParseDuration(*expiration)

	if err != nil {
		panic(err)
	}
	
	exp := time.Now().Add(dur)

	fmt.Println(buck.SignedURL(*key, exp))
}
