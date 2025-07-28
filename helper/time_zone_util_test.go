package helper

import "testing"

func TestConvertUTCtoTimeZone(t *testing.T) {
	tests := []struct {
		name      string
		utcTime   string
		timeZone  string
		want      string
		wantError bool
	}{
		{
			name:      "UTC to PST",
			utcTime:   "2022/01/01 00:00:00",
			timeZone:  "America/Los_Angeles",
			want:      "2021/12/31 16:00:00",
			wantError: false,
		},
		{
			name:      "Invalid Timezone",
			utcTime:   "2022/01/01 00:00:00",
			timeZone:  "Invalid/Timezone",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertUTCtoTimeZone(tt.utcTime, tt.timeZone)
			if (err != nil) != tt.wantError {
				t.Errorf("ConvertUTCtoTimeZone() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if !tt.wantError && got != tt.want {
				t.Errorf("ConvertUTCtoTimeZone() = %v, want %v", got, tt.want)
			}
		})
	}
}
