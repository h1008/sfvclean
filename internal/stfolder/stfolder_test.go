package stfolder

import (
	"testing"
)

func Test_extractFilename(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantTime string
		wantErr  bool
	}{
		{
			name: "With Extension",
			args: args{
				file: "some/path/to/file~20200711-205951.txt",
			},
			wantName: "some/path/to/file.txt",
			wantTime: "2020-07-11 20:59:51 +0000 UTC",
			wantErr:  false,
		},
		{
			name: "Without Extension",
			args: args{
				file: "some/path/to/file~20200711-205951",
			},
			wantName: "some/path/to/file",
			wantTime: "2020-07-11 20:59:51 +0000 UTC",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotTime, err := extractFilename(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("extractFilename() gotName = %v, wantName %v", gotName, tt.wantName)
			}
			if gotTime.String() != tt.wantTime {
				t.Errorf("extractFilename() gotTime = %v, wantTime %v", gotTime, tt.wantTime)
			}
		})
	}
}
