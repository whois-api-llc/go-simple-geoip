package simplegeoip

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathGeoipResponseOK         = "/Geoip/ok"
	pathGeoipResponseError      = "/Geoip/error"
	pathGeoipResponse500        = "/Geoip/500"
	pathGeoipResponsePartial1   = "/Geoip/partial"
	pathGeoipResponsePartial2   = "/Geoip/partial2"
	pathGeoipResponseUnparsable = "/Geoip/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the IP Geolocation API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathGeoipResponseOK:
		case pathGeoipResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathGeoipResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathGeoipResponsePartial1:
			response = response[:len(response)-10]
		case pathGeoipResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathGeoipResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new IP Geolocation API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:   apiServer.Client(),
		GeoipBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestGeoipGet tests the Get function.
func TestGeoipGet(t *testing.T) {
	checkResultRec := func(res *GeoIPResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"ip":"8.8.8.8","location":{"country":"US","region":"California","city":"Mountain View",
"lat":37.38605,"lng":-122.08385,"postalCode":"94035","timezone":"-07:00","geonameId":5375480},
"domains":["000000-1v1v1v1v1v1v118888888.sdqpwlbock-gkynimr.tokyo"],
"as":{"asn":15169,"name":"GOOGLE","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/","type":"Content"},
"isp":"Google LLC","connectionType":""}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"error":"test error message"}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type args struct {
		ctx     context.Context
		options string
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathGeoipResponseOK,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathGeoipResponse500,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathGeoipResponsePartial1,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathGeoipResponsePartial2,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathGeoipResponseError,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			want:    false,
			wantErr: "API error: [499] test error message",
		},
		{
			name: "unparsable response",
			path: pathGeoipResponseUnparsable,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.Get(tt.args.ctx)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("Geoip.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("Geoip.Get() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("Geoip.Get() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestGeoipAPIGetRaw tests the RawData function.
func TestGeoipGetRaw(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"ip":"8.8.8.8","location":{"country":"US","region":"California","city":"Mountain View",
"lat":37.38605,"lng":-122.08385,"postalCode":"94035","timezone":"-07:00","geonameId":5375480},
"domains":["000000-1v1v1v1v1v1v118888888.sdqpwlbock-gkynimr.tokyo"],
"as":{"asn":15169,"name":"GOOGLE","route":"8.8.8.0\/24","domain":"https:\/\/about.google\/intl\/en\/","type":"Content"},
"isp":"Google LLC","connectionType":""}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"code":499,"error":"test error message"}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type args struct {
		ctx     context.Context
		options string
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathGeoipResponseOK,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathGeoipResponse500,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathGeoipResponsePartial1,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathGeoipResponsePartial2,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathGeoipResponseUnparsable,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathGeoipResponseError,
			args: args{
				ctx:     ctx,
				options: "8.8.8.8",
			},
			wantErr: "API failed with status code: 499",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRaw(tt.args.ctx)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("Geoip.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !checkResultRaw(resp.Body) {
				t.Errorf("Geoip.Get() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
