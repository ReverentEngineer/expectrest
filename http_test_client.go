package main

/*
#include <rdg.h>
#include <stdlib.h>
#cgo LDFLAGS: -lrdg
*/
import "C"

import (
	"net/http"
	"strings"
	"sync"
)

type response struct {
	Code    int               `json:"code"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type HTTPTestResult struct {
	Success bool   `json:"success"`
	Err     string `json:"err,omitempty"`

	Request HTTPTestSpec `json:"testSpecification"`

	Response *response `json:"response,omitempty"`
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClientFactory interface {
	New() (HTTPClient, error)
}
type HTTPTestClient struct {
	ClientFactory HTTPClientFactory
}

func (c HTTPTestClient) RunTest(test HTTPTestSpec, ch chan HTTPTestResult, wg *sync.WaitGroup) {
	var req *http.Request
	var err error

	result := HTTPTestResult{
		Request: test,
	}

	if test.Body != nil {
		req, err = http.NewRequest(test.Method, test.Url, strings.NewReader(*test.Body))
	} else {
		req, err = http.NewRequest(test.Method, test.Url, nil)
	}

	if err != nil {
		result.Success = false
		result.Err = err.Error()
		ch <- result
		return
	}

	client, err := c.ClientFactory.New()
	resp, err := client.Do(req)

	if err != nil {
		result.Success = false
		result.Err = err.Error()
		ch <- result
		return
	}

	defer resp.Body.Close()

	result.Response = &response{
		Code: resp.StatusCode,
	}
	result.Success = true
	ch <- result
	return
}
