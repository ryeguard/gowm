package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ryeguard/gowm/internal"
	"github.com/ryeguard/gowm/onecall"
	"github.com/ryeguard/gowm/pkg/owm"
)

var httpAddr = flag.String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
var owmAppID string

func init() {
	flag.StringVar(&owmAppID, "api-key", "", "OpenWeatherMap API key")
}

type weatherClient struct {
	client *owm.Client
}

type GetWeatherArgs struct {
	Location string `json:"location" mcp:"the place's name to get weather for, on the format 'city,country'" jsonschema:"the place's name to get weather for, on the format 'city,country'"`
}

type GetWeatherResult struct {
	Data *onecall.CurrentResponse
}

func (w *weatherClient) GetWeather(ctx context.Context, req *mcp.CallToolRequest, args *GetWeatherArgs) (*mcp.CallToolResult, *GetWeatherResult, error) {
	response, err := w.client.GetWeather(args.Location, &onecall.OneCallOptions{
		Exclude: []onecall.Part{onecall.Parts.MINUTELY, onecall.Parts.HOURLY, onecall.Parts.ALERTS},
	})
	if err != nil {
		return nil, nil, fmt.Errorf("get weather: %w", err)
	}

	weather := response.OneCall
	content := fmt.Sprintf("Weather for %s:\n\n", args.Location)

	content += fmt.Sprintf("Current Temperature: %.1f°C\n", weather.Current.Temp)
	content += fmt.Sprintf("Feels like: %.1f°C\n", weather.Current.FeelsLike)
	content += fmt.Sprintf("Humidity: %d%%\n", weather.Current.Humidity)
	content += fmt.Sprintf("Pressure: %d hPa\n", weather.Current.Pressure)

	if len(weather.Current.Weather) > 0 {
		content += fmt.Sprintf("Conditions: %s\n", weather.Current.Weather[0].Description)
	}

	content += fmt.Sprintf("\nDaily forecast (%d days):\n", len(weather.Daily))
	for i, day := range weather.Daily {
		if i >= 7 { // Limit to next 7 days
			break
		}
		content += fmt.Sprintf("Day %d: %.1f°C / %.1f°C", i+1, day.Temp.Min, day.Temp.Max)
		if len(day.Weather) > 0 {
			content += fmt.Sprintf(" - %s", day.Weather[0].Description)
		}
		content += "\n"
	}

	return &mcp.CallToolResult{Content: []mcp.Content{&mcp.TextContent{Text: content}}},
		&GetWeatherResult{Data: &response.OneCall.Current},
		nil
}

func main() {
	flag.Parse()

	var opts owm.ClientOptions
	if owmAppID != "" {
		opts.AppID = owmAppID
	} else if appID, ok := internal.LoadEnvVar(); ok {
		opts.AppID = appID
	} else {
		log.Fatalln("OpenWeatherMap API key must be set as environment variable or command line flag")
	}

	wc := weatherClient{client: owm.NewClient(&opts).WithOneCall(nil).WithGeo(nil)}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "gowm-api",
		Version: "0.1.0",
		Title:   "OpenWeatherMap weather data",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_weather",
		Description: "Get the full weather forecast for a location (city,country). Always provide only the location here, not the date or time.",
	}, wc.GetWeather)

	if *httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("MCP handler listening at %s", *httpAddr)
		http.ListenAndServe(*httpAddr, handler)
	} else {
		log.Printf("MCP running on stdio")
		t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
		if err := server.Run(context.Background(), t); err != nil {
			log.Printf("Server failed: %v", err)
		}
	}
}
