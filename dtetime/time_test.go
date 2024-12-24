package dtetime_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/peterHoburg/go-date-and-time-extension/dtetime"
)

func ExampleDate() {
	dteTime, err := dtetime.Parse("15:04:05")
	if err != nil {
		panic(">" + err.Error() + "<")
	}

	marshaled, err := json.Marshal(dteTime)
	if err != nil {
		panic(">" + err.Error() + "<")
	}
	fmt.Println(string(marshaled))

	unmarshalled := dtetime.Time{}
	err = json.Unmarshal(marshaled, &unmarshalled)
	if err != nil {
		panic(">" + err.Error() + "<")
	}
	fmt.Println(unmarshalled)

	// Output: "15:04:05"
	// 15:04:05
}

func TestMarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		inputTime dtetime.Time
		want      string
		wantError bool
	}{
		{
			name:      "valid time",
			inputTime: func() dtetime.Time { dteTime, _ := dtetime.Parse("15:04:05"); return dteTime }(),
			want:      `"15:04:05"`,
			wantError: false,
		},
		{
			name:      "zero time",
			inputTime: dtetime.Time{},
			want:      `"00:00:00"`,
			wantError: false,
		},
		{
			name:      "invalid time format (manually created object)",
			inputTime: dtetime.Time{Time: time.Date(0, 0, 0, -1, 0, 0, 0, time.UTC)},
			want:      `"23:00:00"`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(tt.inputTime)
			if (err != nil) != tt.wantError {
				t.Errorf("MarshalJSON() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if string(got) != tt.want {
				t.Errorf("MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
