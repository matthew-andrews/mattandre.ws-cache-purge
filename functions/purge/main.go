package main

import (
	"encoding/json"
	"fmt"
	"github.com/apex/go-apex"
	"time"
)

type message struct {
	China string `json:"china"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		london, _ := time.LoadLocation("Europe/London")
		china := time.Date(2016, 8, 9, 0, 0, 0, 0, london)
		now := time.Now()
		totalDays := int(china.Sub(now).Hours() / 24)
		weeks := (totalDays - totalDays%7) / 7
		days := totalDays % 7
		output := fmt.Sprintf("%d weeks, %d days", weeks, days)
		return message{output}, nil
	})
}
