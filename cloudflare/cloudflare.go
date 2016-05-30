package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type PurgeCacheResponseBody struct {
	Result struct {
		Id string `json:"id"`
	} `json:"result"`
	Success bool `json:"success"`
}

type PurgeCacheRequestBody struct {
	Files []string `json:"files"`
}

type Client struct {
	BaseURL    string
	HttpClient *http.Client
}

func (m *Client) Purge(urls []string) (PurgeCacheResponseBody, error) {
	var responseBody PurgeCacheResponseBody

	// Build and make request
	url := m.BaseURL + "/client/v4/zones/" + os.Getenv("CLOUDFLARE_IDENTIFIER") + "/purge_cache"
	requestBody, err := json.Marshal(PurgeCacheRequestBody{urls})
	if err != nil {
		fmt.Fprintf(os.Stderr, "event=UNEXPECTED_ERROR at=MARSHAL_REQUEST_BODY error=%s\n", err)
		return responseBody, err
	}
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Fprintf(os.Stderr, "event=UNEXPECTED_ERROR at=CREATE_REQUEST error=%s\n", err)
		return responseBody, err
	}
	req.Header.Add("X-Auth-Email", os.Getenv("CLOUDFLARE_AUTH_EMAIL"))
	req.Header.Add("X-Auth-Key", os.Getenv("CLOUDFLARE_API_KEY"))
	req.Header.Add("Content-Type", "application/json")
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "event=UNEXPECTED_ERROR at=DO_REQUEST error=%s\n", err)
		return responseBody, err
	}

	// Open and unmarshal response
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "event=UNEXPECTED_ERROR at=READ_BODY error=%s\n", err)
		return responseBody, err
	}
	if err := json.Unmarshal(bytes, &responseBody); err != nil {
		fmt.Fprintf(os.Stderr, "event=UNEXPECTED_ERROR at=UNMARSHAL_RESPONSE_BODY json=%s error=%s\n", bytes, err)
		return responseBody, err
	}

	return responseBody, nil
}
