package main

import (
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/s3"
	"fmt"
	"flag"
	"time"
	"os"
	"io"
	"net/http"
	"crypto/tls"
)

func main() {
	bucket := flag.String("bucket", "", "bucket")
	key := flag.String("key", "", "key")
	region := flag.String("region", aws.USEast.Name, "region")
	expiration := flag.String("expiration", "1h", "expiration in the Go duration format.  A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as \"300ms\", \"-1.5h\" or \"2h45m\". Valid time units are \"ns\", \"us\" (or \"µs\"), \"ms\", \"s\", \"m\", \"h\".")
	verbose := flag.Bool("v", false, "verbose")
	mode := flag.String("mode", "sign", "sign or cat")
	insecure := flag.Bool("insecure", false, "turn on InsecureSkipVerify: http://golang.org/pkg/crypto/tls/#Config")

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

	if *verbose {
		fmt.Fprintf(os.Stderr, "auth: %#v\n", auth)
	}

	reg, found := aws.Regions[*region]

	if ! found {
		panic(fmt.Sprintf("invalid region: %s", *region))
	}

	s3c := s3.New(auth, reg)
	buck := s3c.Bucket(*bucket)

	switch *mode {
	case "sign":
		dur, err := time.ParseDuration(*expiration)

		if err != nil {
			panic(err)
		}
		
		exp := time.Now().Add(dur)

		fmt.Println(buck.SignedURL(*key, exp))

	case "cat":
		if *insecure {
			tr := &http.Transport{
				TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			}
			http.DefaultClient.Transport = tr
		}
		
		reader, err := buck.GetReader(*key)

		if err != nil {
			panic(err)
		}

		io.Copy(os.Stdout, reader)
		
	default:
		panic(fmt.Sprintf("invalid mode: %s", *mode))
	}
}
