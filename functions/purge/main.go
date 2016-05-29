package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/apex/go-apex"
	"github.com/matthew-andrews/mattandre.ws-websitecdnpurge/s3eventtourls"
	"io/ioutil"
	"net/http"
	"os"
)

type PurgeCacheRequestBody struct {
	Files []string `json:"files"`
}

type PurgeCacheResponseBody struct {
	Result struct {
		Id string `json:"id"`
	} `json:"result"`
	Success bool `json:"success"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		urls, err := s3eventtourls.S3EventToUrls(event)
		if err != nil {
			return nil, err
		}

		fmt.Fprintf(os.Stderr, "WILL PURGE: %s\n", urls)
		requestBody, err := json.Marshal(PurgeCacheRequestBody{
			Files: urls,
		})
		if err != nil {
			return nil, err
		}

		body := bytes.NewBuffer(requestBody)

		client := &http.Client{}
		url := "https://api.cloudflare.com/client/v4/zones/" + os.Getenv("CLOUDFLARE_IDENTIFIER") + "/purge_cache"
		req, err := http.NewRequest("DELETE", url, body)
		if err != nil {
			return nil, err
		}

		req.Header.Add("X-Auth-Email", os.Getenv("CLOUDFLARE_AUTH_EMAIL"))
		req.Header.Add("X-Auth-Key", os.Getenv("CLOUDFLARE_API_KEY"))
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var responseBody PurgeCacheResponseBody
		if err := json.Unmarshal(bytes, &responseBody); err != nil {
			return nil, err
		}

		return responseBody, nil
	})
}
