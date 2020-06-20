package pb

import "testing"

func Test_oldVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantOld bool
	}{
		{"protoc-2.5.0", args{version: "2.5.0"}, true},
		{"protoc-2.6.0", args{version: "2.6.0"}, true},
		{"protoc-2.7.0", args{version: "2.7.0"}, true},
		{"protoc-3.5.0", args{version: "3.5.0"}, true},
		{"protoc-3.6.0", args{version: "3.6.0"}, false},
		{"protoc-3.6.1", args{version: "3.6.1"}, false},
		{"protoc-3.7.0", args{version: "3.7.0"}, false},
		{"protoc-3.7.1", args{version: "3.7.1"}, false},
		{"protoc-3.10.1", args{version: "3.10.1"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOld, err := oldVersion(tt.args.version)
			if err != nil {
				t.Errorf("judge protoc version error = %v", err)
			}
			if gotOld != tt.wantOld {
				t.Errorf("oldVersion() gotOld = %v, want %v", gotOld, tt.wantOld)
			}
		})
	}
}
