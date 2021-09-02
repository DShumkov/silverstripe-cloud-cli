package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

//FetchResponse contains response from API
type FetchResponse struct {
	Data struct {
		Type       string `json:"type"`
		ID         int    `json:"id"`
		Attributes struct {
			Status string      `json:"status"`
			Error  interface{} `json:"error"`
		} `json:"attributes"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
	Meta struct {
		ServerTime string `json:"server_time"`
		StatusCode int    `json:"status_code"`
	} `json:"meta"`
}

func (fr *FetchResponse) hydrate(jsonString []byte) error {
	return json.Unmarshal(jsonString, fr)
}

//FetchStatusResponse contains response from API
type FetchStatusResponse struct {
	Data struct {
		Type       string `json:"type"`
		ID         int    `json:"id"`
		Attributes struct {
			Status string `json:"status"`
		} `json:"attributes"`
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
	} `json:"data"`
}

func (fsr *FetchStatusResponse) hydrate(jsonString []byte) error {
	return json.Unmarshal(jsonString, fsr)
}

//RunFetch adding new fetch
func RunFetch(c *Config, projectID string) {
	url := strings.Trim(c.BaseURL, "/") + "/project/" + projectID + "/git/fetches"
	Log.Debug("Fetching: " + url)
	response, err := HttpCall(c.Email, c.APIToken, url, "POST", nil)
	check(err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	check(err)

	var fr FetchResponse

	err = fr.hydrate(body)
	check(err)

	if fr.Data.ID > 0 {
		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			if CheckFetch(c, &fr, projectID) {
				Log.Debug("Fetch completed")
				return
			}
		}
		Log.Warn("Failed fetching. Timeout!")
		return
	}
	Log.Err("Error fetch response. ID is invalid")
}

//CheckFetch checking fetch status
func CheckFetch(c *Config, fr *FetchResponse, projectID string) bool {
	if 0 == fr.Data.ID {
		Log.Err("Invalid fetching process ID")
	}
	url := strings.Trim(c.BaseURL, "/") + "/project/" + projectID + "/git/fetches/" + strconv.Itoa(fr.Data.ID)
	response, err := HttpCall(c.Email, c.APIToken, url, "GET", nil)
	check(err)

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	check(err)
	var result FetchStatusResponse
	json.Unmarshal(body, &result)
	if result.Data.Attributes.Status == "Complete" {
		return true
	}
	return false
}
