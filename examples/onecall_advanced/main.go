package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ryeguard/gowm/onecall"
	"golang.org/x/time/rate"
)

func main() {
	httpClient := &http.Client{
		Timeout: time.Second,

		// We can implement client-side rate limiting
		// by defining a custom transport.
		Transport: newRateLimitedTransport(rate.Every(time.Second), 1, nil),
	}

	client, err := onecall.NewClient(&onecall.ClientOptions{
		HttpClient: httpClient,
		AppID:      "YOUR-API-KEY",

		// By default, OpenWeatherMap API returns Kelvin for temperature,
		// which is not very common for everyday applications.
		Units: onecall.Units.METRIC,
	})
	if err != nil {
		panic(err)
	}
	resp, err := client.OneCall(59.3327, 18.0656, &onecall.OneCallOptions{
		// If we only want CURRENT and DAILY for our location we can exclude the other forecasts.
		Exclude: []onecall.Exclude{onecall.Excludes.HOURLY, onecall.Excludes.MINUTELY, onecall.Excludes.ALERTS},

		// Setting `Units` on the OneCall call will overrule the one set on the client.
		Units: onecall.Units.IMPERIAL,

		// Some string fields from OpenWeatherMap are in the local language.
		Lang: onecall.Langs.SWEDISH,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("The temperature at %v is %v but feels like %v\n", resp.Current.Dt.Format(time.Kitchen), resp.Current.Temp, resp.Current.FeelsLike)

	now := time.Now()
	var overmorrow onecall.DailyResponse
	overmorrowFound := false
	for _, d := range resp.Daily {
		if d.Dt.YearDay() == now.YearDay()+1 {
			overmorrow = d
			overmorrowFound = true
			break
		}
	}
	if overmorrowFound {
		fmt.Printf("The weather the day after tomorrow will be '%v' with a max temp of %v\n", overmorrow.Summary, overmorrow.Temp.Max)
	} else {
		fmt.Println("The forecast did not include the weather for the day after tomorrow.")
	}
}

type rateLimitedTransport struct {
	transport http.RoundTripper
	limiter   *rate.Limiter
}

func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.transport.RoundTrip(req)
}

func newRateLimitedTransport(limit rate.Limit, burst int, transport http.RoundTripper) *rateLimitedTransport {
	if transport == nil {
		transport = http.DefaultTransport
	}

	return &rateLimitedTransport{
		transport: transport,
		limiter:   rate.NewLimiter(limit, burst),
	}
}
