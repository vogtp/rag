package vecdb

import (
	"testing"
)

// 2025-09-11 16:30:36.875593874 +0200 CEST m=+207.045403821
// 2006-01-02 15:04:05.999999999 -0700
// 2025-10-23 10:29:07.4234687 +0200 CEST m=+133.822645301
func Test_parseTime(t *testing.T) {
	tests := []struct {
		timeStr string
		y       int
		m       int
		d       int
		h       int
		min     int
		s       int
		wantErr bool
	}{
		{timeStr: "2025-09-11 16:30:36.875593874 +0200 CEST m=+207.045403821", y: 2025, m: 9, d: 11, h: 16, min: 30, s: 36, wantErr: false},
		{timeStr: "2006-01-02 15:04:05.999999999 -0700", y: 2006, m: 1, d: 2, h: 15, min: 04, s: 5, wantErr: false},
		{timeStr: "2025-10-23 10:29:07.4234687 +0200 CEST m=+133.822645301", y: 2025, m: 10, d: 23, h: 10, min: 29, s: 7, wantErr: false},
		{timeStr: "2025--23 10:29:07.4234687 +0200 CEST m=+133.822645301", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.timeStr, func(t *testing.T) {
			got, gotErr := parseTime(tt.timeStr)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parseTime() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("parseTime() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if got.Year() != tt.y || int(got.Month()) != tt.m || got.Day() != tt.d || got.Hour() != tt.h || got.Minute() != tt.min || got.Second() != tt.s {
				t.Errorf("parseTime() = %s", got)
			}
		})
	}
}
