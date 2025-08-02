package main

import (
	"fmt"
	"time"

	"github.com/ryeguard/gowm/onecall"
)

func main() {
	client, err := onecall.NewClient(&onecall.ClientOptions{
		AppID: "YOUR-API-KEY",

		// By default, OpenWeatherMap API returns Kelvin for temperature,
		// which is not very common for everyday applications.
		Units: onecall.Units.METRIC,
	})
	if err != nil {
		panic(err)
	}

	resp, err := client.OneCall(59.3327, 18.0656, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("The temperature at %v is %v but feels like %v", resp.Current.Dt.Format(time.Kitchen), resp.Current.Temp, resp.Current.FeelsLike)
}
