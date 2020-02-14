package utils

import (
	"net/http"
	"net/http/cookiejar"
	"strings"
)

// Client - http client definition
type Client struct {
	httpclient *http.Client
	Header     *http.Header
	jar        *cookiejar.Jar
}

// CreateClient - creates a http client
func CreateClient() *Client {
	client := Client{}
	client.jar, _ = cookiejar.New(nil)
	client.httpclient = &http.Client{
		Jar: client.jar,
	}
	return &client
}

// Perform - executes the request
func (client *Client) Perform(method, url string, data *strings.Reader) *http.Response {
	var request *http.Request
	var err error
	if data == nil {
		request, err = http.NewRequest(method, url, nil)
	} else {
		request, err = http.NewRequest(method, url, data)
	}
	if err != nil {
		panic(err)
	}
	// request.Header = client.Headers
	// request.Header = *client.Headers.Clone()
	request.Header = CloneHeader(*client.Header)
	response, err := client.httpclient.Do(request)
	if err != nil {
		panic(err)
	}
	return response
}

// CloneHeader - creates a copy of the client headers
func CloneHeader(in http.Header) http.Header {
	out := make(http.Header, len(in))
	for key, values := range in {
		newValues := make([]string, len(values))
		copy(newValues, values)
		out[key] = newValues
	}
	return out
}
