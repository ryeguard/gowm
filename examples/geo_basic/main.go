package main

import (
	"fmt"
	"log"

	"github.com/ryeguard/gowm/geo"

	_ "github.com/joho/godotenv/autoload" // auto-loads .env file
)

func main() {
	client, err := geo.NewClient(nil)
	if err != nil {
		log.Fatalf("new client: %v", err)
	}

	resp, err := client.Direct("Stockholm,SE", nil)
	if err != nil {
		log.Fatalf("direct: %v", err)
	}

	fmt.Printf("%v places matched the query (max 5)\n", len(resp.Data))

	fmt.Printf("%v, %v, or '%v' is located at %v,%v", resp.Data[0].Name, resp.Data[0].Country, resp.Data[0].LocalNames["es"], resp.Data[0].Lat, resp.Data[0].Lon)
}
