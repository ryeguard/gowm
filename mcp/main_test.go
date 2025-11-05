package main

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ryeguard/gowm/pkg/owm"
	"github.com/stretchr/testify/require"
)

// TestWeatherClientInitialization tests that weatherClient can be created
func TestWeatherClientInitialization(t *testing.T) {
	opts := &owm.ClientOptions{
		AppID: "test-api-key",
	}

	wc := weatherClient{
		client: owm.NewClient(opts).WithOneCall(nil).WithGeo(nil),
	}

	require.NotNil(t, wc.client)
}

// TestMCPServerInitialization tests that MCP server can be initialized
func TestMCPServerInitialization(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "gowm-api",
		Version: "0.1.0",
		Title:   "OpenWeatherMap weather data",
	}, nil)

	require.NotNil(t, server)
}

// TestAddToolToServer tests that a tool can be added to the server
func TestAddToolToServer(t *testing.T) {
	opts := &owm.ClientOptions{
		AppID: "test-api-key",
	}

	wc := weatherClient{
		client: owm.NewClient(opts).WithOneCall(nil).WithGeo(nil),
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "gowm-api",
		Version: "0.1.0",
		Title:   "OpenWeatherMap weather data",
	}, nil)

	// This should not panic
	require.NotPanics(t, func() {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "get_weather",
			Description: "Get the full weather forecast for a location (city,country). Always provide only the location here, not the date or time.",
		}, wc.GetWeather)
	})
}

// TestGetWeatherArgs tests the GetWeatherArgs structure
func TestGetWeatherArgs(t *testing.T) {
	args := &GetWeatherArgs{
		Location: "Stockholm,Sweden",
	}

	require.NotNil(t, args)
	require.Equal(t, "Stockholm,Sweden", args.Location)
}

// TestGetWeatherResult tests the GetWeatherResult structure
func TestGetWeatherResult(t *testing.T) {
	result := &GetWeatherResult{
		Data: nil,
	}

	require.NotNil(t, result)
	require.Nil(t, result.Data)
}

// TestGetWeatherWithInvalidLocation tests GetWeather with context
// This tests that the method signature is correct and doesn't panic on setup
func TestGetWeatherWithInvalidLocation(t *testing.T) {
	opts := &owm.ClientOptions{
		AppID: "test-api-key",
	}

	wc := weatherClient{
		client: owm.NewClient(opts).WithOneCall(nil).WithGeo(nil),
	}

	ctx := context.Background()
	req := &mcp.CallToolRequest{}
	args := &GetWeatherArgs{
		Location: "TestCity,TestCountry",
	}

	// This will likely fail with an API error, but we're just testing
	// that the method can be called and returns the expected types
	result, weatherResult, err := wc.GetWeather(ctx, req, args)

	// We expect an error since we're using a fake API key and location
	// The important thing is that the method signature works
	_ = result
	_ = weatherResult
	_ = err

	// This test passes if it doesn't panic - the method is callable
}
