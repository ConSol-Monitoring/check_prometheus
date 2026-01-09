package helper

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// TimestampFreshness is the amount of second a result is treated as valid
var TimestampFreshness int

// InsecureSkipVerify will skip TLS certificate verification when set to true. It will be used when constructing http Transport
var InsecureSkipVerify bool

// Cookies parsed into []*http.Cookie
var Cookies []*http.Cookie

// Verbose flag writes to here
var Verbose bool

type prometheusInterceptor struct {
	next http.RoundTripper
}

// Interceptor function used in verbose mode
func (i *prometheusInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	if Verbose {
		fmt.Printf("Sending %s request to %s\n", req.Method, req.URL.String())
		fmt.Printf("Request:\n%+v\n", req)
		fmt.Printf("Url:\n%+v\n", req.URL)
		fmt.Printf("Header:\n%+v\n", req.Header)

		// Read and print the body content
		if req.Body != nil {
			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				fmt.Printf("Error reading body: %v\n", err)
			} else {
				fmt.Printf("Body:\n%s\n", string(bodyBytes))
				// Restore the body for further processing
				req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
			}
		} else {
			fmt.Printf("Body is empty\n")
		}
	}

	// 2. Ensure the Content-Type is set
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return i.next.RoundTrip(req)
}

// NewAPIClientV1 will create an prometheus api client v1
func NewAPIClientV1(address *url.URL) (v1.API, error) {
	baseTransport := http.DefaultTransport.(*http.Transport).Clone()
	baseTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: InsecureSkipVerify}

	interceptedTransport := &prometheusInterceptor{
		next: baseTransport,
	}

	httpClient := &http.Client{
		Transport: interceptedTransport,
	}

	// Initialize cookie jar only when Cookies are provided
	if len(Cookies) > 0 {
		jar, _ := cookiejar.New(nil)
		httpClient.Jar = jar
		httpClient.Jar.SetCookies(address, Cookies)
	}

	prometheusClient, err := api.NewClient(api.Config{
		Address: address.String(),
		Client:  httpClient,
	})

	if err != nil {
		return nil, err
	}

	return v1.NewAPI(prometheusClient), nil
}

// DoAPIRequest does the http handling for an api request
func DoAPIRequest(url *url.URL) ([]byte, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: InsecureSkipVerify}

	httpClient := &http.Client{
		Transport: transport,
	}

	if len(Cookies) > 0 {
		jar, _ := cookiejar.New(nil)
		httpClient.Jar = jar
		httpClient.Jar.SetCookies(url, Cookies)
	}

	resp, err := httpClient.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// CheckTimestampFreshness tests if the data is still valid
func CheckTimestampFreshness(timestamp model.Time) error {
	return CheckTimeFreshness(time.Unix(int64(timestamp), 0))
}

// CheckTimeFreshness tests if the data is still valid
func CheckTimeFreshness(timestamp time.Time) error {
	if TimestampFreshness == 0 {
		return fmt.Errorf("error when checking time freshness, timestampFreshness is zero")
	}
	timeDiff := time.Since(timestamp)
	if int(timeDiff.Seconds()) > TimestampFreshness {
		return fmt.Errorf("one of the scraped data exceed the freshness by %ds", int(timeDiff.Seconds())-TimestampFreshness)
	}
	return nil
}
