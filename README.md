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

### Geocoding API

The [Geocoding API](https://openweathermap.org/api/geocoding-api) client is implemented in [`geo/client.go`](./geo/client.go). The available methods are:

- `Direct`: Allows you to get geographical coordinates (lat, lon) by using name of the location (city name or area name).
- `Reverse`: Allows you to get name of the location (city name or area name) by using geographical coordinates (lat, lon).

## Features

### MCP Server

This repo implements a LLM-friendly MCP server for the OpenWeatherMap APIs. To use the server, first build the binary:

```bash
go build -o bin/mcp ./mcp
```

Then, configure your client, e.g., Claude Desktop, to use the binary. As of writing (August 2025), this is done in the `claude_desktop_config.json` found by navigating to Settings > Developer > Edit Config in the Claude Desktop application.

```json
{
  "mcpServers": {
    "weather": {
      "command": "PATH/TO/REPO/gowm/bin/mcp",
      "args": [],
      "env": {
        "OWM_API_KEY": "YOUR_API_KEY"
      }
    }
  }
}
```

where `PATH/TO/REPO/gowm/bin/mcp` is the absolute path to the binary and `YOUR_API_KEY` is the OpenWeatherMap API key you can get from signing up/logging in at [openweathermap.org](https://openweathermap.org/).

### CLI

This repo implements a CLI for interacting with the OpenWeatherMap APIs. It may be used as follows:

```bash
# With an existing Go installation:
go run ./cmd/... get-weather 'stockholm,sweden' --api-key=YOUR_API_KEY

# Or, after installation of the CLI binary (instructions to be added):
gowm get-weather 'stockholm,sweden' --api-key=YOUR_API_KEY
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
