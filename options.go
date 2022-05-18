package simplegeoip

import (
	"net/url"
	"strconv"
	"strings"
)

// Option adds parameters to the query.
type Option func(v url.Values)

var _ = []Option{
	OptionOutputFormat("JSON"),
	OptionIPAddress("8.8.8.8"),
	OptionDomain("whoisxmlapi.com"),
	OptionEmail("support@whoisxmlapi.com"),
	OptionReverseIP(0),
}

// OptionOutputFormat sets Response output format JSON | XML. Default: JSON.
func OptionOutputFormat(outputFormat string) Option {
	return func(v url.Values) {
		v.Set("outputFormat", strings.ToUpper(outputFormat))
	}
}

// OptionIPAddress sets IPv4 or IPv6 to search location by.
// If the parameter is not specified, then it defaults to client request's public IP address.
func OptionIPAddress(value string) Option {
	return func(v url.Values) {
		v.Set("ipAddress", value)
	}
}

// OptionDomain sets the domain name to search location by.
// If the parameter is not specified, then 'ipAddress' will be used.
func OptionDomain(value string) Option {
	return func(v url.Values) {
		v.Set("domain", value)
	}
}

// OptionEmail sets the email address or domain name to search location by its MX servers.
// If the parameter is not specified, then 'domain' will be used.
func OptionEmail(value string) Option {
	return func(v url.Values) {
		v.Set("email", value)
	}
}

// OptionReverseIP sets the parameter for showing 5 domains associated with the IP address. Default: 1.
func OptionReverseIP(value int) Option {
	return func(v url.Values) {
		v.Set("reverseIp", strconv.Itoa(value))
	}
}
