package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ryeguard/gowm/onecall"
	"github.com/ryeguard/gowm/pkg/owm"

	_ "github.com/joho/godotenv/autoload" // auto-loads .env file
)

func main() {
	client := owm.
		NewClient(nil).
		WithOneCall(&onecall.ClientOptions{
			// By default, OpenWeatherMap API returns Kelvin for temperature,
			// which is not very common for everyday applications.
			Units: onecall.Units.METRIC,
		}).
		WithGeo(nil)

	geo, err := client.Geo.Direct("Stockholm,SE", nil)
	if err != nil {
		log.Fatalf("geo Direct: %v", err)
	}

	fmt.Printf("%v (%v) is located at %.4f,%.4f\n", geo.Data[0].Name, geo.Data[0].Country, geo.Data[0].Lat, geo.Data[0].Lon)

	oc, err := client.OneCall.OneCall(59.3327, 18.0656, nil)
	if err != nil {
		log.Fatalf("onecall OneCall: %v", err)
	}

	fmt.Printf("The temperature at %v is %v but feels like %v\n", oc.Current.Dt.Format(time.Kitchen), oc.Current.Temp, oc.Current.FeelsLike)
}
