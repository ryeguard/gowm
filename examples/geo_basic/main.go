package main

import (
	"fmt"
	"log"

	"github.com/ryeguard/gowm/geo"

	_ "github.com/joho/godotenv/autoload" // auto-loads .env file
)

func main() {
	client := geo.NewClient(nil)

	// Comma-separated name (city) and ISO 3166 country code
	q := "Stockholm,SE"
	directResp, err := client.Direct(q, nil)
	if err != nil {
		log.Fatalf("direct: %v", err)
	}
	if len(directResp.Data) == 0 {
		log.Fatalf("no results for %v", q)
	}

	fmt.Printf("%v places matched the query (max 5)\n", len(directResp.Data))

	stockholm := directResp.Data[0]

	fmt.Printf("%v (%v), or '%v' in Spanish, is located at %.4f, %.4f\n", stockholm.Name, stockholm.Country, stockholm.LocalNames["es"], stockholm.Lat, stockholm.Lon)

	reverseResp, err := client.Reverse(stockholm.Lat-1, stockholm.Lon-1, nil)
	if err != nil {
		log.Fatalf("direct: %v", err)
	}

	fmt.Printf("South-west of %v we find %v\n", stockholm.Name, reverseResp.Data[0].Name)
}
