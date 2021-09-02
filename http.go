package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func HttpCall(username string, password string, url, method string, body io.Reader) (*http.Response, error) {
	client := &http.Client{
		Timeout: time.Second * 60,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("X-Api-Version", "2.0")

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}

	return response, nil
}
