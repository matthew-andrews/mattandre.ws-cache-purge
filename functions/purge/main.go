package main

import (
	"encoding/json"
	"errors"
	"github.com/apex/go-apex"
	"net/http"
	"os"
	"strings"
)

type message struct {
	Url string `json:"url"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var m message

		if err := json.Unmarshal(event, &m); err != nil {
			return nil, err
		}

		body := strings.NewReader("{\"files\":[\"" + m.Url + "/\"]}")

		client := &http.Client{}
		if req, err := http.NewRequest("DELETE", "https://api.cloudflare.com/client/v4/"+os.Getenv("CLOUDFLARE_IDENTIFIER")+"/purge_cache", body); err != nil {
			return nil, err
		}

		req.Header.Add("X-Auth-Email", "matt@mattandre.ws")
		req.Header.Add("X-Auth-Key", os.Getenv("CLOUDFLARE_API_KEY"))
		req.Header.Add("Content-Type", "application/json")

		if resp, err := client.Do(req); err != nil {
			return nil, err
		}
		if resp.StatusCode > 399 {
			return errors.New("Bad Server Response: " + resp.Status)
		}
		return m, nil
	})
}
