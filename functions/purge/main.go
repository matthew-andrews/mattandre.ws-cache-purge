package main

import (
	"encoding/json"
	"fmt"
	"github.com/apex/go-apex"
	"github.com/matthew-andrews/mattandre.ws-websitecdnpurge/cloudflare"
	"github.com/matthew-andrews/mattandre.ws-websitecdnpurge/s3eventtourls"
	"net/http"
	"os"
)

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {

		// Get URLs from S3 event
		urls, err := s3eventtourls.S3EventToUrls(event)
		if err != nil {
			fmt.Fprintf(os.Stderr, "event=UNEXPECTED_ERROR error=%s\n", err)
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "event=WILL_PURGE urls=%s\n", urls)

		// Purge those URLs from CloudFlare
		cloudFlareClient := &cloudflare.Client{"https://api.cloudflare.com", &http.Client{}}
		cloudFlareResponse, err := cloudFlareClient.Purge(urls)
		if err != nil {
			fmt.Fprintf(os.Stderr, "event=UNEXPECTED_ERROR error=%s\n", err)
			return nil, err
		}
		return cloudFlareResponse, nil
	})
}
