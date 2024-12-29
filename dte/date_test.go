package dte_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/peterHoburg/go-date-and-time-extension/dte"
)

func ExampleNewDate() {
	dteTime, err := dte.NewDate("2006-01-02")
	if err != nil {
		return
	}

	fmt.Println(dteTime)

	// Output: 2006-01-02
}

func ExampleDate_json_to_struct() {
	type TestStruct struct {
		Date dte.Date `json:"date"`
	}

	testStruct := TestStruct{}

	err := json.Unmarshal([]byte(`{"date":"2006-01-02"}`), &testStruct)
	if err != nil {
		return
	}

	fmt.Println(testStruct.Date)

	// Output: 2006-01-02
}

func ExampleDate_struct_to_json() {
	type TestStruct struct {
		Date dte.Date `json:"date"`
	}

	testStruct := TestStruct{}

	parsed, err := dte.NewDate("2006-01-02")
	if err != nil {
		return
	}

	testStruct.Date = parsed

	marshaled, err := json.Marshal(testStruct)
	if err != nil {
		return
	}

	fmt.Println(string(marshaled))

	// Output: {"date":"2006-01-02"}
}

//nolint:funlen
func TestDateNewDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		inputTime string
		want      string
		wantError bool
	}{
		{
			name:      "invalid date",
			inputTime: "0-01-01",
			want:      ``,
			wantError: true,
		},
		{
			name:      "invalid date TZ",
			inputTime: "2006-01-02Z",
			want:      ``,
			wantError: true,
		},
		{
			name:      "valid date",
			inputTime: "2006-01-02",
			want:      `"2006-01-02"`,
			wantError: false,
		},
		{
			name:      "zero time",
			inputTime: "0000-01-01",
			want:      `"0000-01-01"`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			parsed, err := dte.NewDate(tt.inputTime)
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

func TestDateSetFromTime(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		inputTime time.Time
		want      string
		wantError bool
	}{
		{
			name:      "valid",
			inputTime: time.Date(2023, 10, 15, 20, 4, 5, 0, time.FixedZone("UTC+5", 5*3600)),
			want:      `"2023-10-15"`,
			wantError: false,
		},
		{
			name:      "valid time with +20 TZ",
			inputTime: time.Date(2023, 10, 15, 10, 4, 5, 0, time.FixedZone("UTC+20", 20*3600)),
			want:      `"2023-10-15"`,
			wantError: false,
		},
		{
			name:      "valid time UTC",
			inputTime: time.Date(2023, 10, 15, 15, 4, 5, 0, time.UTC),
			want:      `"2023-10-15"`,
			wantError: false,
		},
		{
			name:      "invalid empty time",
			inputTime: time.Time{},
			want:      `"0001-01-01"`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var dteDate dte.Date

			err := dteDate.SetFromTime(tt.inputTime)
			if (err != nil) != tt.wantError {
				t.Errorf("SetFromTime() error = %v, wantError %v", err, tt.wantError)

				return
			}

			got, err := json.Marshal(dteDate)
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

//nolint:funlen
func TestDateUnmarshalJSON(t *testing.T) { //nolint:dupl
	t.Parallel()

	type TestStruct struct {
		Date dte.Date `json:"date"`
	}

	tests := []struct {
		name      string
		inputJSON string
		want      dte.Date
		wantErr   bool
	}{
		{
			name:      "valid date",
			inputJSON: "{\n\t\"date\": \"2006-01-02\"\n}",

			want:    func() dte.Date { t, _ := dte.NewDate("2006-01-02"); return t }(), //nolint:nlreturn
			wantErr: false,
		},
		{
			name:      "null date",
			inputJSON: "{\n\t\"date\": \"null\"\n}",

			want:    func() dte.Date { t, _ := dte.NewDate("null"); return t }(), //nolint:nlreturn
			wantErr: false,
		},
		{
			name:      "int",
			inputJSON: "{\n\t\"date\": 1\n}",

			want:    func() dte.Date { t, _ := dte.NewDate("null"); return t }(), //nolint:nlreturn
			wantErr: true,
		},
		{
			name:      "invalid json date",
			inputJSON: "{\n\t\"date\": 2006-01-02\n}",

			want:    func() dte.Date { t, _ := dte.NewDate("null"); return t }(), //nolint:nlreturn
			wantErr: true,
		},
		{
			name:      "Full timestamp",
			inputJSON: "{\n\t\"date\": \"2006-01-02T15:04:05Z\"\n}",
			want:      func() dte.Date { t, _ := dte.NewDate("2006-01-02"); return t }(), //nolint:nlreturn
			wantErr:   false,
		},
		{
			name:      "Full timestamp with -20",
			inputJSON: "{\n\t\"date\": \"2006-01-02T10:04:05-20:00\"\n}",
			want:      func() dte.Date { t, _ := dte.NewDate("2006-01-02"); return t }(), //nolint:nlreturn
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var testStruct TestStruct

			err := json.Unmarshal([]byte(tt.inputJSON), &testStruct)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if err != nil {
				return
			}

			if testStruct.Date != tt.want {
				t.Errorf("UnmarshalJSON() = %v, want %v", testStruct.Date, tt.want)

				return
			}
		})
	}
}
