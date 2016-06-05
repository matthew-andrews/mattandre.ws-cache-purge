package s3eventtourls

import (
	"encoding/json"
	"reflect"
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

	if !reflect.DeepEqual([]string{"https://mattandre.ws/cv/index.html", "https://mattandre.ws/cv/", "https://mattandre.ws/cv"}, actual) {
		t.Fatalf("%s does not match expected value", actual)
	}
}
