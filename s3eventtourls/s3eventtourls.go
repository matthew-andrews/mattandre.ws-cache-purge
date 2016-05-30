package s3eventtourls

import (
	"encoding/json"
	"regexp"
)

type S3Update struct {
	Records []struct {
		S3 struct {
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
		} `json:"s3"`
	} `json:"Records"`
}

func S3EventToUrls(event json.RawMessage) ([]string, error) {
	var updates S3Update

	if err := json.Unmarshal(event, &updates); err != nil {
		return nil, err
	}

	var urls []string
	for _, update := range updates.Records {
		url := "https://" + update.S3.Bucket.Name + "/" + update.S3.Object.Key
		urls = append(urls, url)

		// If ends with index.html also purge that URL minus index.html
		re := regexp.MustCompile("/index\\.html$")
		if trimmedUrl := re.ReplaceAllString(url, ""); trimmedUrl != url {
			urls = append(urls, trimmedUrl+"/")
			urls = append(urls, trimmedUrl)
		}
	}

	return urls, nil
}
