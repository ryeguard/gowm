# gowm

![Go Version](https://img.shields.io/github/go-mod/go-version/ryeguard/gowm)
![License](https://img.shields.io/github/license/ryeguard/gowm)
![Build Status](https://img.shields.io/github/actions/workflow/status/ryeguard/gowm/test.yml?branch=main)

A Go [OpenWeatherMap API](https://openweathermap.org/api) client. The latest Go docs may be browsed at [pkg.go.dev](https://pkg.go.dev/github.com/ryeguard/gowm).

## Getting Started

First, make sure you have all of the below prerequisites:

- Installation of Go, e.g., from [go.dev](https://go.dev/doc/install).
- API key from OpenWeatherMap
  - Sign up for free at [openweathermap.org](https://openweathermap.org/)
  - Get your API key after having signed in
  - Note: API key activation may take some time

Now, you should be able to either use this project's packages in your project, or clone this repo and contribute/run the examples provided.

Examples of basic as well as advanced usage are located in the [examples](./examples/) dir. To run any of the examples (having cloned this repo):

```bash
# e.g., the owm_basic example
go run examples/owm_basic/main.go
```

## APIs

### One Call API 3.0

The [One Call API 3.0](https://openweathermap.org/api/one-call-3) client is implemented in [`onecall/client.go`](./onecall/client.go). The available methods are:

- `CurrentAndForecast`: To get access to current weather, minute forecast for 1 hour, hourly forecast for 48 hours, daily forecast for 8 days and government weather alerts.
  - `CurrentAndForecastRaw` is available and returns a direct mapping of the API response, e.g., not parsing timestamps to `time.Time` but rather leaving them as `int`s.

- `Historical`: To get historical weather data for a specific timestamp. Historical data is available from January 1, 1979 onwards.
  - `HistoricalRaw` is available and returns a direct mapping of the API response.

### Geocoding API

The [Geocoding API](https://openweathermap.org/api/geocoding-api) client is implemented in [`geo/client.go`](./geo/client.go). The available methods are:

- `Direct`: Allows you to get geographical coordinates (lat, lon) by using name of the location (city name or area name).
- `Reverse`: Allows you to get name of the location (city name or area name) by using geographical coordinates (lat, lon).

## Features

### MCP Server

This repo implements a LLM-friendly MCP (Model Context Protocol) server for the OpenWeatherMap APIs. The MCP server is built into the main `gowm` binary as a subcommand.

The MCP server provides the following tools:
- `get_weather`: Get current weather and forecast for a location
- `get_historical_weather`: Get historical weather data for a specific location and date

#### Usage

Start the MCP server using the `gowm mcp` command:

```bash
# Start MCP server in stdio mode (default, for Claude Desktop)
gowm mcp --api-key=YOUR_API_KEY

# Or set API key via environment variable
export OWM_API_KEY=YOUR_API_KEY
gowm mcp

# Start in HTTP mode for testing
gowm mcp --http=localhost:8080
```

#### Configuration for Claude Desktop

Configure your client, e.g., Claude Desktop, to use the MCP server. As of writing (January 2025), this is done in the `claude_desktop_config.json` found by navigating to Settings > Developer > Edit Config in the Claude Desktop application.

```json
{
  "mcpServers": {
    "weather": {
      "command": "/usr/local/bin/gowm",
      "args": ["mcp"],
      "env": {
        "OWM_API_KEY": "YOUR_API_KEY"
      }
    }
  }
}
```

where `/usr/local/bin/gowm` is the absolute path to the binary (adjust based on your installation location - use `which gowm` to find it) and `YOUR_API_KEY` is the OpenWeatherMap API key you can get from signing up/logging in at [openweathermap.org](https://openweathermap.org/).

## Installation

The `gowm` binary includes both CLI and MCP server functionality in a single unified tool.

**Option 1: Quick install script (macOS/Linux)**

```bash
curl -sSL https://raw.githubusercontent.com/ryeguard/gowm/main/install.sh | bash
```

**Option 2: Download pre-built binary**

Download the latest `gowm` binary for your platform from the [releases page](https://github.com/ryeguard/gowm/releases/latest).

**Option 3: Install with Go**

```bash
go install github.com/ryeguard/gowm/cmd/gowm@latest
```

**Option 4: Build from source**

```bash
git clone https://github.com/ryeguard/gowm.git
cd gowm
go build -o bin/gowm ./cmd/gowm
```

### CLI

Use the CLI to interact with OpenWeatherMap APIs from the command line.

#### Usage

```bash
# Get weather forecast for a location
gowm get-weather 'stockholm,sweden' --api-key=YOUR_API_KEY

# Get historical weather for a location and date
gowm get-historical-weather 'stockholm,sweden' '2024-01-15' --api-key=YOUR_API_KEY

# Or set the API key via environment variable
export OWM_API_KEY=YOUR_API_KEY
gowm get-weather 'stockholm,sweden'
gowm get-historical-weather 'new york,us' '2024-01-15T12:00:00Z'

# Check version
gowm version
```

### Static Types

Leveraging Go's type system, as well as generating better go enums using [`zarldev/goenums`](https://github.com/zarldev/goenums), using the clients is straight-forward. You don't need to worry about guessing the input format of the API calls, of for example languages and units. Rather than:

```go
// from briandowns/openweathermap (another great OpenWeatherMap Go client and the inspiration for this project)
w, err := owm.NewOneCall("F", "EN", apiKey, []string{})
if err != nil {
  log.Fatalln(err)
}

err = w.OneCallByCoordinates(
  &Coordinates{
    Latitude:  33.45,
    Longitude: -112.07,
  },
)
if err != nil {
  t.Error(err)
}
```

We can instead do:

```go
client := onecall.NewClient(&onecall.ClientOptions{
  AppID: apiKey
  Units: onecall.Units.IMPERIAL,
})

resp, err := client.CurrentAndForecast(33.45, -112.07, &onecall.OneCallOptions{
  Lang: onecall.Langs.ENGLISH,
})
if err != nil {
  log.Fatalln(err)
}
```

### Custom `http.Client`s and `slog.Logger`s

You can pass custom HTTP clients and loggers to the API Client to make the most of Go's std lib features like rate limiting and structured logging with configurable logging levels.

## Contributing

Contributions are welcome.

## Disclaimer

This library is an unofficial Go client for the OpenWeather API. It is not affiliated with or endorsed by OpenWeather. See the license at [LICENSE](./LICENSE).

Use of this client requires a valid API key from [OpenWeather](https://openweathermap.org/), and use of OpenWeather data is subject to their [license terms](https://openweathermap.org/price). Please ensure you comply with their data licensing conditions, particularly around attribution and share-alike requirements.

## Links

- [ryeguard/gowm Go docs](https://pkg.go.dev/github.com/ryeguard/gowm)
- [OpenWeatherMap API](https://openweathermap.org/api)
- [zarldev/goenums project](https://github.com/zarldev/goenums) used for generating type-safe enums.
