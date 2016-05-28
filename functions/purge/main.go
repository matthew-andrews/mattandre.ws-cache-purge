package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/apex/go-apex"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

type S3Update struct {
	Records []struct {
		S3 struct {
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

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
		var updates S3Update

		if err := json.Unmarshal(event, &updates); err != nil {
			return nil, err
		}

		urls := make([]string, len(updates.Records), 2*len(updates.Records))
		for i, update := range updates.Records {
			url := "https://mattandre.ws/" + update.S3.Object.Key
			fmt.Fprintf(os.Stderr, "WILL PURGE: "+url+"\n")
			urls[i] = url

			// If ends with index.html also purge that URL minus index.html
			re := regexp.MustCompile("index\\.html$")
			if trimmedUrl := re.ReplaceAllString(url, ""); trimmedUrl != url {
				fmt.Fprintf(os.Stderr, "WILL PURGE: "+trimmedUrl+"\n")
				urls[i] = trimmedUrl
			}
		}

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

		req.Header.Add("X-Auth-Email", "matt@mattandre.ws")
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
