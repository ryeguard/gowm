package main

import (
	"net/http"

	"github.com/ryeguard/gowm/onecall"
	"github.com/ryeguard/gowm/pkg/owm"
)

func main() {
	client := owm.NewClient(
		http.DefaultClient,
		"SECRET",
		&onecall.ClientOptions{
			Units: onecall.Units.METRIC,
		})

	_, err := client.OneCall.OneCall(59.3327, 18.0656, &onecall.OneCallOptions{
		Exclude: []onecall.Exclude{onecall.Excludes.CURRENT},
	})
	if err != nil {
		panic(err)
	}
}
