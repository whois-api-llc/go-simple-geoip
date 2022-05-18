package simplegeoip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GeoipService is an interface for IP Geolocation API.
type GeoipService interface {
	// Get returns parsed IP Geolocation API response
	Get(ctx context.Context, opts ...Option) (*GeoIPResponse, *Response, error)

	// GetRaw returns raw IP Geolocation API response as Response struct with Body saved as a byte slice
	GetRaw(ctx context.Context, opts ...Option) (*Response, error)
}

// Response is the http.Response wrapper with Body saved as a byte slice.
type Response struct {
	*http.Response

	// Body is the byte slice representation of http.Response Body
	Body []byte
}

// geoipServiceOp is the type implementing the GeoipService interface.
type geoipServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ GeoipService = &geoipServiceOp{}

// newRequest creates the API request with default parameters and the specified apiKey.
func (service *geoipServiceOp) newRequest() (*http.Request, error) {
	req, err := service.client.NewRequest(http.MethodGet, service.baseURL, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("apiKey", service.client.apiKey)

	req.URL.RawQuery = query.Encode()

	return req, nil
}

// apiResponse is used for parsing IP Geolocation API response as a model instance.
type apiResponse struct {
	GeoIPResponse
	Code    int    `json:"code"`
	Message string `json:"error"`
}

// request returns intermediate API response for further actions.
func (service *geoipServiceOp) request(ctx context.Context, opts ...Option) (*Response, error) {
	req, err := service.newRequest()
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	for _, opt := range opts {
		opt(q)
	}

	req.URL.RawQuery = q.Encode()

	var b bytes.Buffer

	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw IP Geolocation API response.
func parse(raw []byte) (*apiResponse, error) {
	var response apiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// Get returns parsed IP Geolocation API response.
func (service geoipServiceOp) Get(
	ctx context.Context,
	opts ...Option,
) (geoipResponse *GeoIPResponse, resp *Response, err error) {
	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	geoipResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if geoipResp.Message != "" || geoipResp.Code != 0 {
		return nil, nil, &ErrorMessage{
			Code:    geoipResp.Code,
			Message: geoipResp.Message,
		}
	}

	return &geoipResp.GeoIPResponse, resp, nil
}

// GetRaw returns raw IP Geolocation API response as Response struct with Body saved as a byte slice.
func (service geoipServiceOp) GetRaw(
	ctx context.Context,
	opts ...Option,
) (resp *Response, err error) {
	resp, err = service.request(ctx, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error.
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string.
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
