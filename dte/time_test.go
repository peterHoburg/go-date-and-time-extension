package dte_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/peterHoburg/go-date-and-time-extension/dte"
)

func ExampleParse() {
	dteTime, err := dte.NewTime("15:04:05Z")
	if err != nil {
		return
	}
	fmt.Println(dteTime)

	// Output: 15:04:05Z
}

func ExampleTime_json_to_struct() {
	type TestStruct struct {
		Time dte.Time `json:"time"`
	}
	testStruct := TestStruct{}
	err := json.Unmarshal([]byte(`{"time":"15:04:05Z"}`), &testStruct)
	if err != nil {
		return
	}
	fmt.Println(testStruct.Time)

	// Output: 15:04:05Z
}

func ExampleTime_struct_to_json() {
	type TestStruct struct {
		Time dte.Time `json:"time"`
	}

	testStruct := TestStruct{}

	parsed, err := dte.NewTime("15:04:05Z")
	if err != nil {
		return
	}

	testStruct.Time = parsed
	marshaled, err := json.Marshal(testStruct)

	if err != nil {
		return
	}
	fmt.Println(string(marshaled))

	// Output: {"time":"15:04:05Z"}
}

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		inputTime string
		want      string
		wantError bool
	}{
		{
			name:      "invalid time no TZ",
			inputTime: "15:04:05",
			want:      ``,
			wantError: true,
		},
		{
			name:      "invalid time bad TZ",
			inputTime: "15:04:05-55:00:00",
			want:      ``,
			wantError: true,
		},
		{
			name:      "invalid time zulu and TZ",
			inputTime: "15:04:05Z-05:00:00",
			want:      ``,
			wantError: true,
		},
		{
			name:      "invalid time zulu and TZ no -",
			inputTime: "15:04:05Z05:00:00",
			want:      ``,
			wantError: true,
		},
		{
			name:      "valid time zulu",
			inputTime: "15:04:05Z",
			want:      `"15:04:05Z"`,
			wantError: false,
		},
		{
			name:      "valid time -5",
			inputTime: "10:04:05-05:00",
			want:      `"15:04:05Z"`,
			wantError: false,
		},
		{
			name:      "valid time +5",
			inputTime: "20:04:05+05:00",
			want:      `"15:04:05Z"`,
			wantError: false,
		},
		{
			name:      "zero time",
			inputTime: "00:00:00Z",
			want:      `"00:00:00Z"`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := dte.NewTime(tt.inputTime)
			if (err != nil) && true == tt.wantError {
				return
			}
			if err != nil {
				t.Errorf("Parse() error = %v, wantError %v", err, tt.wantError)
				return
			}
			got, err := json.Marshal(parsed)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestSetFromTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		inputTime time.Time
		want      string
		wantError bool
	}{
		{
			name:      "valid specific time",
			inputTime: time.Date(2023, 10, 15, 20, 4, 5, 0, time.FixedZone("UTC+5", 5*3600)),
			want:      `"15:04:05Z"`,
			wantError: false,
		},
		{
			name:      "valid time with negative offset",
			inputTime: time.Date(2023, 12, 25, 10, 4, 5, 0, time.FixedZone("UTC-3", -5*3600)),
			want:      `"15:04:05Z"`,
			wantError: false,
		},
		{
			name:      "valid time UTC",
			inputTime: time.Date(2023, 12, 25, 15, 4, 5, 0, time.UTC),
			want:      `"15:04:05Z"`,
			wantError: false,
		},
		{
			name:      "valid zero time UTC",
			inputTime: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			want:      `"00:00:00Z"`,
			wantError: false,
		},
		{
			name:      "invalid empty time",
			inputTime: time.Time{},
			want:      `"00:00:00Z"`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var dteTime dte.Time
			err := dteTime.SetFromTime(tt.inputTime)
			if (err != nil) != tt.wantError {
				t.Errorf("SetFromTime() error = %v, wantError %v", err, tt.wantError)
				return
			}

			got, err := json.Marshal(dteTime)
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}

			if string(got) != tt.want {
				t.Errorf("SetFromTime() = %v, want %v", string(got), tt.want)
			}
		})
	}
}