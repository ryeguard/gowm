package main

import (
	"fmt"
	"os"

	"github.com/ryeguard/gowm/onecall"
	"github.com/ryeguard/gowm/pkg/owm"
	"github.com/spf13/cobra"
)

var owmClient *owm.Client

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

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().String("api-key", "", "OWM API key")

	rootCmd.AddCommand(getWeatherCmd)
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
