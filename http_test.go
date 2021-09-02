package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpCall(t *testing.T) {
	testUsername := "username"
	testPassword := "12345"
	serverTestMessage := "test response 42"
	clientTestMessage := "I am client 42"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, serverTestMessage)
		u, p, ok := r.BasicAuth()
		if !ok {
			t.Error("Bad basic auth!")
		}
		if testUsername != u {
			t.Error("Bad username:", u)
		}
		if testPassword != p {
			t.Error("Bad password", p)
		}
		body, _ := ioutil.ReadAll(r.Body)
		if strings.TrimSpace(string(body)) != clientTestMessage {
			t.Errorf("Bad request body: %s", string(body))
		}
	}))
	defer ts.Close()

	response, err := HttpCall(testUsername, testPassword, ts.URL, "GET", bytes.NewBufferString(clientTestMessage))
	if nil != err {
		t.Errorf("HTTP call response error: %s", err)
	}
	body, _ := ioutil.ReadAll(response.Body)
	if strings.TrimSpace(string(body)) != serverTestMessage {
		fmt.Println("body:", len(string(body)), "var:", len(serverTestMessage))
		t.Errorf("Bad response body: %s not equal to: %s", string(body), serverTestMessage)
	}

}
