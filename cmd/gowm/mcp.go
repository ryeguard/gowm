package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ryeguard/gowm/internal"
	"github.com/ryeguard/gowm/onecall"
	"github.com/ryeguard/gowm/pkg/owm"
	"github.com/spf13/cobra"
)

// weatherClient wraps the OWM client for MCP tool handlers
type weatherClient struct {
	client *owm.Client
}

// GetWeatherArgs defines the arguments for the get_weather MCP tool
type GetWeatherArgs struct {
	Location string `json:"location" mcp:"the place's name to get weather for, on the format 'city,country'" jsonschema:"the place's name to get weather for, on the format 'city,country'"`
}

// GetWeatherResult defines the result structure for the get_weather MCP tool
type GetWeatherResult struct {
	Data *onecall.CurrentResponse
}

// GetWeather is the MCP tool handler for getting weather data
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

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start the MCP server for LLM integration",
	Long: `Start the Model Context Protocol (MCP) server for integration with LLM clients like Claude Desktop.

The server provides weather data tools that can be used by LLMs to access OpenWeatherMap data.

Examples:
  # Start MCP server in stdio mode (default, for Claude Desktop)
  gowm mcp

  # Start MCP server in HTTP mode
  gowm mcp --http=localhost:8080

Configuration for Claude Desktop (claude_desktop_config.json):
  {
    "mcpServers": {
      "weather": {
        "command": "/path/to/gowm",
        "args": ["mcp"],
        "env": {
          "OWM_API_KEY": "YOUR_API_KEY"
        }
      }
    }
  }`,
	RunE: runMCPServer,
}

func init() {
	mcpCmd.Flags().String("http", "", "if set, use streamable HTTP at this address, instead of stdin/stdout")
}

func runMCPServer(cmd *cobra.Command, args []string) error {
	// Get API key - try flag first, then environment variables
	apiKey, err := cmd.Flags().GetString("api-key")
	if err != nil {
		return err
	}

	var opts owm.ClientOptions
	if apiKey != "" {
		opts.AppID = apiKey
	} else if appID, ok := internal.LoadEnvVar(); ok {
		opts.AppID = appID
	} else {
		return fmt.Errorf("OpenWeatherMap API key must be set via --api-key flag or OWM_API_KEY/OWM_APP_ID environment variable")
	}

	// Create weather client
	wc := weatherClient{client: owm.NewClient(&opts).WithOneCall(nil).WithGeo(nil)}

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "gowm-api",
		Version: version, // Use the version from root.go
		Title:   "OpenWeatherMap weather data",
	}, nil)

	// Add weather tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_weather",
		Description: "Get the full weather forecast for a location (city,country). Always provide only the location here, not the date or time.",
	}, wc.GetWeather)

	// Start server in appropriate mode
	httpAddr, err := cmd.Flags().GetString("http")
	if err != nil {
		return err
	}

	if httpAddr != "" {
		handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
			return server
		}, nil)
		log.Printf("MCP server listening at %s", httpAddr)
		return http.ListenAndServe(httpAddr, handler)
	} else {
		log.Printf("MCP server running on stdio")
		t := &mcp.LoggingTransport{Transport: &mcp.StdioTransport{}, Writer: os.Stderr}
		return server.Run(context.Background(), t)
	}
}
