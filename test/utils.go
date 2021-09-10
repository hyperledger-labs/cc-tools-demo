package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

func GetAndVerify(url string, expectedStatus int, expectedResponse interface{}) error {
	// Disable TLS cert verification
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("wrong status. expected: %d received: %d", expectedStatus, resp.StatusCode)
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var receivedResponse interface{}
	err = json.Unmarshal(resBody, &receivedResponse)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(expectedResponse, receivedResponse) {
		return fmt.Errorf("unexpected response\n\nexpected: %#v\n\nreceived: %#v", expectedResponse, receivedResponse)
	}
	return nil
}

func PostAndVerify(url string, body interface{}, expectedStatus int, expectedResponse interface{}) error {
	postBody, _ := json.Marshal(body)
	postBodyBuf := bytes.NewBuffer(postBody)

	// Disable TLS cert verification
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}

	resp, err := client.Post(url, "application/json", postBodyBuf)
	if err != nil {
		return err
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != expectedStatus {
		return fmt.Errorf("wrong status. expected: %d received: %d\n%s", expectedStatus, resp.StatusCode, string(resBody))
	}
	var receivedResponse interface{}
	err = json.Unmarshal(resBody, &receivedResponse)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(expectedResponse, receivedResponse) {
		return fmt.Errorf("unexpected response\n\nexpected: %#v\n\nreceived: %#v", expectedResponse, receivedResponse)
	}
	return nil
}

func fail() {
	fmt.Println("FAIL")
}

func pass() {
	fmt.Println("PASS")
}
