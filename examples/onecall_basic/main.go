package main

import (
	"github.com/ryeguard/gowm/onecall"
	"github.com/ryeguard/gowm/pkg/owm"
)

func main() {
	client, err := owm.NewClient(
		&onecall.ClientOptions{
			AppID: "SECRET",
			Units: onecall.Units.METRIC,
		})
	if err != nil {
		panic(err)
	}

	_, err = client.OneCall.OneCall(59.3327, 18.0656, &onecall.OneCallOptions{
		Exclude: []onecall.Exclude{onecall.Excludes.CURRENT},
	})
	if err != nil {
		panic(err)
	}
}
