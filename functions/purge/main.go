package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/apex/go-apex"
	"net/http"
	"os"
)

type message struct {
	Urls []string `json:"urls"`
}

type PurgeCacheRequest struct {
	Files []string `json:"files"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		requestObject, err := json.Marshal(PurgeCacheRequest{
			Files: m.Urls,
		})
		if err != nil {
			return nil, err
		}

		body := bytes.NewBuffer(requestObject)

		client := &http.Client{}
		url := "https://api.cloudflare.com/client/v4/zones/" + os.Getenv("CLOUDFLARE_IDENTIFIER") + "/purge_cache"
		fmt.Fprintf(os.Stderr, url)
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
		if resp.StatusCode > 399 {
			return nil, errors.New("Bad Server Response: " + resp.Status)
		}
		return m, nil
	})
}
