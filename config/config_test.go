package config

import (
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		consumerKey    string
		consumerSecret string
		callbackURL    string
		port           string
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "should return a valid value when all args are provided",
			args: args{"consumerKey", "consumerSecret", "callbackURL", "80"},
			want: &Config{"consumerKey", "consumerSecret", "callbackURL", "80"},
		},
		{
			name:    "should return an error when a parameter value is missing",
			args:    args{"consumerKey", "consumerSecret", "callbackURL", ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.consumerKey, tt.args.consumerSecret, tt.args.callbackURL, tt.args.port)
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
