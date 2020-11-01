package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

type DefaultHTTPClientFactory struct{}

func (f DefaultHTTPClientFactory) New() (HTTPClient, error) {
	return &http.Client{}, nil
}

func main() {
	var tests []HTTPTestSpec
	var results []HTTPTestResult

	results = make([]HTTPTestResult, 0)

	if len(os.Args) != 2 {
		log.Print("Usage: expectrest <config>")
		log.Fatal("Invalid arguments")
	}

	configFile := os.Args[1]

	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	if err = json.Unmarshal([]byte(config), &tests); err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}
	resultChannel := make(chan HTTPTestResult)
	testClient := HTTPTestClient{
		ClientFactory: DefaultHTTPClientFactory{},
	}

	for _, test := range tests {
		wg.Add(1)
		go testClient.RunTest(test, resultChannel, &wg)
	}

	go func() {
		for {
			result, ok := <-resultChannel
			if !ok {
				return
			}
			results = append(results, result)
			wg.Done()
		}
	}()

	wg.Wait()
	close(resultChannel)

	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(results)
	return
}
