package owm

import (
	"fmt"

	"github.com/ryeguard/gowm/onecall"
)

type Client struct {
	OneCall *onecall.Client
}

func NewClient(opts *onecall.ClientOptions) (*Client, error) {
	oc, err := onecall.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("new OneCall client: %w", err)
	}
	return &Client{
		OneCall: oc,
	}, nil
}
