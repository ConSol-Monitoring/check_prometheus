package helper

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/consol-monitoring/check_x"
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

type prometheusInterceptor struct {
	next http.RoundTripper
}

func (i *prometheusInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	// 1. You can log the body here to see exactly what is being sent
	fmt.Printf("Sending %s request to %s\n", req.Method, req.URL.String())

	// 2. Ensure the Content-Type is definitely set
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 3. Fix potential Idempotency-Key issues by removing it if it's nil
	// if val, ok := req.Header["Idempotency-Key"]; ok && val == nil {
	// 	delete(req.Header, "Idempotency-Key")
	// }

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
func CheckTimestampFreshness(timestamp model.Time) {
	CheckTimeFreshness(time.Unix(int64(timestamp), 0))
}

// CheckTimeFreshness tests if the data is still valid
func CheckTimeFreshness(timestamp time.Time) {
	if TimestampFreshness == 0 {
		return
	}
	timeDiff := time.Since(timestamp)
	if int(timeDiff.Seconds()) > TimestampFreshness {
		check_x.Exit(check_x.Unknown, fmt.Sprintf("One of the scraped data exceed the freshness by %ds", int(timeDiff.Seconds())-TimestampFreshness))
	}
}
