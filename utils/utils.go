package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
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
		Jar:       client.jar,
		Transport: createTransport(nil),
	}
	return &client
}

// Perform - executes the request
func (client *Client) Perform(method, url string, data *bytes.Buffer) *http.Response {
	var request *http.Request
	var err error
	if data == nil {
		request, err = http.NewRequest(method, url, nil)
	} else {
		request, err = http.NewRequest(method, url, data)
	}
	if err != nil {
		log.Fatal().Err(err)
	}
	// request.Header = client.Headers
	// request.Header = *client.Headers.Clone()
	request.Header = CloneHeader(*client.Header)
	response, err := client.httpclient.Do(request)
	if err != nil {
		log.Fatal().Err(err)
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

// cfbm
func createTransport(proxy *url.URL) *http.Transport {
	transport := &http.Transport{
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          10,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		TLSClientConfig:       createTLSConfig(),
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
	}

	if proxy != nil {
		transport.Proxy = http.ProxyURL(proxy)
	}

	return transport
}

func createTLSConfig() *tls.Config {
	return &tls.Config{
		Rand:                     rand.Reader,
		KeyLogWriter:             nil,
		InsecureSkipVerify:       false,
		PreferServerCipherSuites: true,
	}
}

// CheckContentType - checks for the response content type
func CheckContentType(responseHeader http.Header) {
	if !strings.HasPrefix(responseHeader.Get("Content-Type"), "application/json") {
		log.Warn().Msg("The response is not a JSON")
	}
}
