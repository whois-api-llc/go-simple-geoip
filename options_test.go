package simplegeoip

import (
	"net/url"
	"reflect"
	"testing"
)

// TestOptions tests the Options functions.
func TestOptions(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
		option Option
		want   string
	}{
		{
			name:   "output format",
			values: url.Values{},
			option: OptionOutputFormat("JSON"),
			want:   "outputFormat=JSON",
		},
		{
			name:   "IP address",
			values: url.Values{},
			option: OptionIPAddress("8.8.8.8"),
			want:   "ipAddress=8.8.8.8",
		},
		{
			name:   "domain",
			values: url.Values{},
			option: OptionDomain("whoisxmlapi.com"),
			want:   "domain=whoisxmlapi.com",
		},
		{
			name:   "email",
			values: url.Values{},
			option: OptionEmail("support@whoisxmlapi.com"),
			want:   "email=support%40whoisxmlapi.com",
		},
		{
			name:   "reverse IP",
			values: url.Values{},
			option: OptionReverseIP(0),
			want:   "reverseIp=0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.values)
			if got := tt.values.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
