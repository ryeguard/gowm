package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ryeguard/gowm/onecall"
	"github.com/ryeguard/gowm/pkg/owm"
	"github.com/spf13/cobra"
)

var (
	owmClient *owm.Client
	// Version information - set by GoReleaser via ldflags
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gowm",
	Short: "A CLI for accessing OpenWeatherMap APIs and data",
	Long: `A CLI for accessing OpenWeatherMap APIs and data. For example:
	
	gowm get-weather 'stockholm,sweden'`,
	PersistentPreRunE: setupClient,
}

var getWeatherCmd = &cobra.Command{
	Use:   "get-weather [city,country]",
	Short: "Get weather for a place",
	Long: `Get the weather for a place.

Examples:
  gowm get-weather 'stockholm,sweden`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		resp, err := owmClient.GetWeather(name, &onecall.OneCallOptions{SaveAsJson: fmt.Sprintf("%v.json", name)})
		if err != nil {
			fmt.Println("get weather: v", err)
			return err
		}

		fmt.Printf("Got weather for %v, %v\n", resp.GeoDirect.Name, resp.GeoDirect.Country)
		return nil
	},
}

var getHistoricalWeatherCmd = &cobra.Command{
	Use:   "get-historical-weather [city,country] [date]",
	Short: "Get historical weather for a place and date",
	Long: `Get historical weather data for a specific location and date.
The date can be in RFC3339 format (e.g. '2024-01-15T12:00:00Z') or simple date format (e.g. '2024-01-15').
Historical data is available from January 1, 1979 onwards.

Examples:
  gowm get-historical-weather 'stockholm,sweden' '2024-01-15'
  gowm get-historical-weather 'new york,us' '2024-01-15T12:00:00Z'`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		location := args[0]
		dateStr := args[1]

		// Parse the date
		var dt time.Time
		var err error

		// Try parsing as RFC3339 first
		dt, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			// Try parsing as just a date (YYYY-MM-DD)
			dt, err = time.Parse("2006-01-02", dateStr)
			if err != nil {
				return fmt.Errorf("invalid date format: %s (expected RFC3339 like '2024-01-15T12:00:00Z' or date like '2024-01-15')", dateStr)
			}
			// If just a date was provided, set time to noon UTC
			dt = time.Date(dt.Year(), dt.Month(), dt.Day(), 12, 0, 0, 0, time.UTC)
		}

		// Get coordinates for the location
		geoData, err := owmClient.GetCoordinates(location)
		if err != nil {
			return fmt.Errorf("get coordinates: %w", err)
		}

		if len(geoData) == 0 {
			return fmt.Errorf("location not found: %s", location)
		}

		geo := geoData[0]

		// Get historical weather
		resp, err := owmClient.GetHistoricalWeather(geo.Lat, geo.Lon, dt, &onecall.OneCallOptions{
			Units: onecall.Units.METRIC,
		})
		if err != nil {
			return fmt.Errorf("get historical weather: %w", err)
		}

		// Print the results
		fmt.Printf("Historical weather for %s, %s on %s:\n\n", geo.Name, geo.Country, dt.Format("2006-01-02 15:04:05 MST"))
		fmt.Printf("Temperature: %.1f°C (feels like %.1f°C)\n", resp.Current.Temp, resp.Current.FeelsLike)
		if len(resp.Current.Weather) > 0 {
			fmt.Printf("Conditions: %s\n", resp.Current.Weather[0].Description)
		}
		fmt.Printf("Humidity: %d%%\n", resp.Current.Humidity)
		fmt.Printf("Pressure: %d hPa\n", resp.Current.Pressure)
		fmt.Printf("Wind Speed: %.1f m/s\n", resp.Current.WindSpeed)
		fmt.Printf("Cloudiness: %d%%\n", resp.Current.Clouds)
		fmt.Printf("UV Index: %.1f\n", resp.Current.UVI)

		if !resp.Current.Sunrise.IsZero() {
			fmt.Printf("Sunrise: %s\n", resp.Current.Sunrise.Format("15:04:05 MST"))
		}
		if !resp.Current.Sunset.IsZero() {
			fmt.Printf("Sunset: %s\n", resp.Current.Sunset.Format("15:04:05 MST"))
		}

		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gowm version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", date)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("api-key", "", "OWM API key")

	rootCmd.AddCommand(getWeatherCmd)
	rootCmd.AddCommand(getHistoricalWeatherCmd)
	rootCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(versionCmd)
}

func setupClient(cmd *cobra.Command, args []string) error {
	appID, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}

	owmClient = owm.NewClient(&owm.ClientOptions{
		AppID: appID,
	}).WithOneCall(nil).WithGeo(nil)
	return nil
}
