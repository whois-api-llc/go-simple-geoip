package simplegeoip

import (
	"fmt"
)

// Location is the part of IP Geolocation API response that contains location details.
type Location struct {
	// Country is the two letters country code from ISO 3166.
	Country string `json:"country"`

	// Region is a region.
	Region string `json:"region"`

	// City is a city.
	City string `json:"city"`

	// Lat is a latitude.
	Lat float64 `json:"lat"`

	// Lng is a longitude.
	Lng float64 `json:"lng"`

	// PostalCode is a postal code.
	PostalCode string `json:"postalCode"`

	// Timezone is the timezone in the format "+10:00".
	Timezone string `json:"timezone"`

	// GeonameID is the ID of location in the GeoNames database. The field is omitted if the record is not found.
	GeonameID uint `json:"geonameId"`
}

// AS is an Autonomous System. It works for IPv4 only. The field is omitted if the record is not found.
type AS struct {
	// ASN is the autonomous system number.
	ASN int `json:"asn"`

	// Name is the autonomous system name.
	Name string `json:"name"`

	// Route is the autonomous system route.
	Route string `json:"route"`

	// Domain is the autonomous system website's URL.
	Domain string `json:"domain"`

	// 	Type is the autonomous system type, one of the following: "Cable/DSL/ISP", "Content", "Educational/Research",
	//	"Enterprise", "Non-Profit", "Not Disclosed", "NSP", "Route Server". Empty when unknown.
	Type string `json:"type"`
}

// GeoIPResponse is a response of IP Geolocation API.
type GeoIPResponse struct {
	// IP is an IP address
	IP string `json:"ip"`

	// Location is the part of IP Geolocation API response that contains location details.
	Location Location `json:"location"`

	// ISP is an internet service provider.
	ISP string `json:"isp"`

	// ConnectionType is the connection type which can be one of "modem", "mobile", "broadband", "company".
	ConnectionType string `json:"connectionType"`

	// Domains is the array of domains associated with the IP. The field is omitted if the record is not found.
	// This array is limited to 5 domains.
	Domains []string `json:"domains"`

	// AS is an autonomous system. It works for IPv4 only. The field is omitted if the record is not found.
	AS AS `json:"as"`
}

// ErrorMessage is an error message.
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

// Error returns error message as a string.
func (e ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.Code, e.Message)
}
