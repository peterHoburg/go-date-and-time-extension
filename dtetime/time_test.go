package dtetime_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/peterHoburg/go-date-and-time-extension/dtetime"
)

func ExampleParse() {
	dteTime, err := dtetime.Parse("15:04:05Z")
	if err != nil {
		return
	}
	fmt.Println(dteTime)

	// Output: 15:04:05Z
}

func ExampleTime_json_to_struct() {
	type TestStruct struct {
		Time dtetime.Time `json:"time"`
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
		Time dtetime.Time `json:"time"`
	}

	testStruct := TestStruct{}

	parsed, err := dtetime.Parse("15:04:05Z")
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

			parsed, err := dtetime.Parse(tt.inputTime)
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
