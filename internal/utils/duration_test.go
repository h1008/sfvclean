package utils

import (
	"reflect"
	"testing"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		givenStr string
		want     Duration
		wantErr  bool
	}{
		{name: "Years", givenStr: "42Y", want: Duration{Years: 42}},
		{name: "Months", givenStr: "12M", want: Duration{Months: 12}},
		{name: "Days", givenStr: "7D", want: Duration{Days: 7}},
		{name: "Empty", givenStr: "", wantErr: true},
		{name: "Invalid Number", givenStr: "e3D", wantErr: true},
		{name: "Negative Number", givenStr: "-7D", wantErr: true},
		{name: "Invalid Unit", givenStr: "7G", wantErr: true},
		{name: "Missing Unit", givenStr: "7", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.givenStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
