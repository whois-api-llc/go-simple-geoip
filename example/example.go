package example

import (
	"context"
	"errors"
	simplegeoip "github.com/whois-api-llc/go-simple-geoip"
	"log"
)

func GetData(apikey string) {
	client := simplegeoip.NewBasicClient(apikey)

	// Get parsed IP Geolocation API response as a model instance
	geoipResp, resp, err := client.GeoipService.Get(context.Background(),
		// this option is ignored, as the inner parser works with JSON only
		simplegeoip.OptionOutputFormat("XML"),
		// this option results in searching location by the specified IP address
		// instead of the client request's public IP address
		simplegeoip.OptionIPAddress("8.8.8.8"))

	if err != nil {
		// Handle error message returned by server
		var apiErr *simplegeoip.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	//Some values are not always returned and need to be validated before using
	if geoipResp.Domains != nil {
		log.Printf("IP: %s, domains: %s, lat: %f, lng: %f\n",
			geoipResp.IP,
			geoipResp.Domains,
			geoipResp.Location.Lat,
			geoipResp.Location.Lng)
	}

	log.Println("raw response is always in JSON format. Most likely you don't need it.")
	log.Printf("raw response: %s\n", string(resp.Body))
}

func GetRawData(apikey string) {
	client := simplegeoip.NewBasicClient(apikey)

	// Get raw API response
	resp, err := client.GeoipService.GetRaw(context.Background(),
		simplegeoip.OptionOutputFormat("JSON"),
		// this option causes reverse IP search not to be performed
		simplegeoip.OptionReverseIP(0))

	if err != nil {
		// Handle error message returned by server
		log.Fatal(err)
	}

	log.Println(string(resp.Body))
}
