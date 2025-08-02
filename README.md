# gowm

A Go [OpenWeatherMap API](https://openweathermap.org/api) client.

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

The [One Call API 3.0](https://openweathermap.org/api/one-call-3) client is implemented in [`onecall/client.go`](./onecall/client.go)

## Contributing

Contributions are welcome.

## Disclaimer

This library is an unofficial Go client for the OpenWeather API. It is not affiliated with or endorsed by OpenWeather. See the license at [LICENSE](./LICENSE).

Use of this client requires a valid API key from [OpenWeather](https://openweathermap.org/), and use of OpenWeather data is subject to their [license terms](https://openweathermap.org/price). Please ensure you comply with their data licensing conditions, particularly around attribution and share-alike requirements.
