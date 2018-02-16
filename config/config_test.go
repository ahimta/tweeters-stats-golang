package config

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		consumerKey    string
		consumerSecret string
		callbackURL    string
		port           string
		homepage       string
		host           string
		protocol       string
		corsDomain     string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "should return a valid value when all args are provided",
			args: args{
				"consumerKey",
				"consumerSecret",
				"callbackURL",
				"8",
				"/",
				"d",
				"h",
				"p",
			},
			want: &Config{
				"consumerKey",
				"consumerSecret",
				"callbackURL",
				"8",
				"/",
				"d",
				"h",
				"p",
			},
		},
		{
			name: "should return a valid value when corsDomain is missing",
			args: args{
				"consumerKey",
				"consumerSecret",
				"callbackURL",
				"80",
				"/",
				"h",
				"p",
				"",
			},
			want: &Config{
				"consumerKey",
				"consumerSecret",
				"callbackURL",
				"80",
				"/",
				"h",
				"p",
				"",
			},
		},
		{
			name: "should return an error when a parameter value is missing",
			args: args{
				"consumerKey",
				"consumerSecret",
				"callbackURL",
				"80",
				"/",
				"h",
				"",
				"",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(
				tt.args.consumerKey,
				tt.args.consumerSecret,
				tt.args.callbackURL,
				tt.args.port,
				tt.args.homepage,
				tt.args.host,
				tt.args.protocol,
				tt.args.corsDomain,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
