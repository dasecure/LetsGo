package main

import (
	"testing"
	"time"

	"snippetbox.dasecure.com/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		in   time.Time
		out  string
	}{
		{
			name: "UTC",
			in:   time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			out:  "17 Mar 2022 at 10:15",
		},
		{
			name: "CET",
			in:   time.Date(2023, 4, 18, 11, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			out:  "18 Apr 2023 at 10:15",
		},
		{
			name: "Empty",
			in:   time.Time{},
			out:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.in)
			assert.Equal(t, hd, tt.out)
		})
	}
}

// func TestHumanDate(t *testing.T) {
// 	tm := time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC)
// 	hd := humanDate(tm)
// 	if hd != "17 Mar 2022 at 10:15" {
// 		t.Errorf("Expected '17 Mar 2022 at 10:15', got '%s'", hd)
// 	}
// }
