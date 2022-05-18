[![go-simple-geoip license](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)
[![go-simple-geoip made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://pkg.go.dev/github.com/whois-api-llc/go-simple-geoip)
[![go-simple-geoip test](https://github.com/whois-api-llc/go-simple-geoip/workflows/Test/badge.svg)](https://github.com/whois-api-llc/go-simple-geoip/actions/)

# Overview

The client library for
[IP Geolocation API](https://ip-geolocation.whoisxmlapi.com)
in Go language.

The minimum go version is 1.17.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/go-simple-geoip
```

# Examples

Full API documentation available [here](https://ip-geolocation.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := simplegeoip.NewBasicClient(apiKey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := simplegeoip.NewClient(apiKey, simplegeoip.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

IP Geolocation API lets you check geographical location by IP address, domain name or email address. 

```go

// Make request to get parsed IP Geolocation API response for the client's public IP address
geoipResp, resp, err := client.GeoipService.Get(ctx)
if err != nil {
    log.Fatal(err)
}

log.Println(geoipResp.IP)
log.Println(geoipResp.Location.Country)

// Make request to get raw IP Geolocation API data
resp, err := client.GeoipService.GetRaw(ctx)
if err != nil {
    log.Fatal(err)
}

log.Println(string(resp.Body))


```
