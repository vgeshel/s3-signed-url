package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"fmt"
	"flag"
	"time"
	"os"
	"io"
)

func main() {
	bucket := flag.String("bucket", "", "bucket")
	key := flag.String("key", "", "key")
	region := flag.String("region", "us-east-1", "region")
	timing := flag.Bool("t", false, "show timing")
	mode := flag.String("mode", "sign", "sign or cat")

	flag.Parse()

	if *bucket == "" || *key == "" {
		fmt.Printf("bucket and key are required\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	s3c := s3.New(session.New(&aws.Config{
		Region: aws.String(*region),
	}))
	input := &s3.GetObjectInput{
		Bucket: aws.String(*bucket),
		Key:    aws.String(*key),
	}

	switch *mode {
	case "cat":
		result, err := s3c.GetObject(input)

		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchKey:
					fmt.Println(s3.ErrCodeNoSuchKey, aerr.Error())
				default:
					fmt.Println(aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			panic(err)
		}

		before := time.Now()
		io.Copy(os.Stdout, result.Body)

		if *timing {
			fmt.Fprintf(os.Stderr, "done %s/%s in %v\n", *bucket, *key, time.Now().Sub(before))
		}

	default:
		panic(fmt.Sprintf("invalid mode: %s", *mode))
	}
}
