package dtetime_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/peterHoburg/go-date-and-time-extension/dtetime"
)

func ExampleTime_MarshalJSON() {
	dteTime, err := dtetime.Parse("15:04:05Z")
	if err != nil {
		return
	}

	marshaled, err := json.Marshal(dteTime)
	if err != nil {
		return
	}
	fmt.Println(string(marshaled))

	unmarshalled := dtetime.Time{}
	err = json.Unmarshal(marshaled, &unmarshalled)
	if err != nil {
		return
	}
	fmt.Println(unmarshalled)

	// Output: "15:04:05Z"
	// 15:04:05Z
}

func ExampleTime_JSONToStruct() {
	type TestStruct struct {
		Time dtetime.Time `json:"time"`
	}
	testStruct := TestStruct{}
	err := json.Unmarshal([]byte(`{"time":"15:04:05Z"}`), &testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Println(testStruct.Time)

	// Output: 15:04:05Z
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
			name:      "valid time",
			inputTime: "15:04:05Z",
			want:      `"15:04:05Z"`,
			wantError: false,
		},
		{
			name:      "valid time",
			inputTime: "15:04:05-05:00",
			want:      `"20:04:05Z"`,
			wantError: false,
		},
		{
			name:      "valid time",
			inputTime: "15:04:05+05:00",
			want:      `"10:04:05Z"`,
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
			if (err != nil) == tt.wantError {
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
