package main

import (
	"fmt"
	"time"

	"github.com/ryeguard/gowm/onecall"

	_ "github.com/joho/godotenv/autoload" // auto-loads .env file
)

func main() {
	client := onecall.NewClient(&onecall.ClientOptions{
		Units: onecall.Units.METRIC,
	})

	// Get historical weather for Stockholm, Sweden on January 1, 2024 at noon UTC
	historicalDate := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	fmt.Printf("Fetching historical weather for %s...\n\n", historicalDate.Format("2006-01-02 15:04:05 MST"))

	resp, err := client.Historical(59.3327, 18.0656, historicalDate, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Location: %s (UTC%+d)\n", resp.Timezone, resp.TimezoneOffset/3600)
	fmt.Printf("Coordinates: %.4f°N, %.4f°E\n\n", resp.Lat, resp.Lon)

	fmt.Printf("Weather at %s:\n", resp.Current.Dt.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("  Temperature: %.1f°C (feels like %.1f°C)\n", resp.Current.Temp, resp.Current.FeelsLike)
	fmt.Printf("  Conditions: %s\n", resp.Current.Weather[0].Description)
	fmt.Printf("  Humidity: %d%%\n", resp.Current.Humidity)
	fmt.Printf("  Pressure: %d hPa\n", resp.Current.Pressure)
	fmt.Printf("  Wind Speed: %.1f m/s\n", resp.Current.WindSpeed)
	fmt.Printf("  Cloudiness: %d%%\n", resp.Current.Clouds)
	fmt.Printf("  UV Index: %.1f\n", resp.Current.UVI)

	if !resp.Current.Sunrise.IsZero() {
		fmt.Printf("  Sunrise: %s\n", resp.Current.Sunrise.Format("15:04:05 MST"))
	}
	if !resp.Current.Sunset.IsZero() {
		fmt.Printf("  Sunset: %s\n", resp.Current.Sunset.Format("15:04:05 MST"))
	}

	// Example: Get historical weather for multiple dates
	fmt.Printf("\n--- Temperature trends for the past week ---\n")
	now := time.Now()
	for i := 7; i >= 1; i-- {
		date := now.AddDate(0, 0, -i)
		// Use noon UTC for consistency
		date = time.Date(date.Year(), date.Month(), date.Day(), 12, 0, 0, 0, time.UTC)

		resp, err := client.Historical(59.3327, 18.0656, date, &onecall.OneCallOptions{
			Units: onecall.Units.METRIC,
		})
		if err != nil {
			fmt.Printf("%s: Error - %v\n", date.Format("2006-01-02"), err)
			continue
		}

		fmt.Printf("%s: %.1f°C - %s\n",
			date.Format("2006-01-02"),
			resp.Current.Temp,
			resp.Current.Weather[0].Description,
		)
	}
}
