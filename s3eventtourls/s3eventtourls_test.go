package s3eventtourls

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestS3EventToUrls(t *testing.T) {
	var jsonRaw json.RawMessage
	var jsonBytes = []byte(`{ "Records": [ {
		"s3": {
			"object": { "key": "cv/index.html" },
			"bucket": { "name": "mattandre.ws" }
		}
	} ] }`)
	json.Unmarshal(jsonBytes, &jsonRaw)

	actual, err := S3EventToUrls(jsonRaw)

	if err != nil {
		t.Fatalf("S3EventToUrls errored incorrectly %s", err)
	}

	err = slicesAreEquivalent([]string{"https://mattandre.ws/cv/index.html", "https://mattandre.ws/cv/", "https://mattandre.ws/cv"}, actual)

	if err != nil {
		t.Fatalf(err.Error())
	}

}
func slicesAreEquivalent(expected []string, actual []string) error {
	if len(expected) != len(actual) {
		return errors.New(fmt.Sprintf("returned slice should be length %d", len(expected)))
	}

	for i, _ := range expected {
		if expected[i] != actual[i] {
			return errors.New(fmt.Sprintf("item at index %d should match %s but was %s", i, expected[i], actual[i]))
		}
	}

	return nil
}
