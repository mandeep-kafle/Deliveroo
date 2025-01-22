package cronparser

import (
	"reflect"
	"testing"
)

func TestCronExpressionParser(t *testing.T) {
	tests := []struct {
		name        string
		expression  string
		want        *CronSchedule
		wantErr     bool
		errorString string
	}{
		{
			name:        "Invalid cron - too few fields",
			expression:  "* * * *",
			want:        nil,
			wantErr:     true,
			errorString: "invalid cron expression: requires at least 6 fields",
		},
		{
			name:        "Invalid cron - invalid minute value",
			expression:  "60 * * * * /usr/bin/find",
			want:        nil,
			wantErr:     true,
			errorString: "error in field 1 (60): value must be between 0 and 59",
		},
		{
			name:       "Basic cron - all wildcards",
			expression: "* * * * * /usr/bin/find",
			want: &CronSchedule{
				Minutes:     makeRange(0, 59),
				Hours:       makeRange(0, 23),
				DaysOfMonth: makeRange(1, 31),
				Months:      makeRange(1, 12),
				DaysOfWeek:  makeRange(0, 6),
				Command:     "/usr/bin/find",
			},
			wantErr: false,
		},
		{
			name:       "Basic cron - specific values",
			expression: "0 0 1 1 0 /usr/bin/find",
			want: &CronSchedule{
				Minutes:     []int{0},
				Hours:       []int{0},
				DaysOfMonth: []int{1},
				Months:      []int{1},
				DaysOfWeek:  []int{0},
				Command:     "/usr/bin/find",
			},
			wantErr: false,
		},
		{
			name:       "Complex cron - range with step",
			expression: "1-7/2 0 1,31 * 1-5 /usr/bin/find",
			want: &CronSchedule{
				Minutes:     []int{1, 3, 5, 7},
				Hours:       []int{0},
				DaysOfMonth: []int{1, 31},
				Months:      makeRange(1, 12),
				DaysOfWeek:  []int{1, 2, 3, 4, 5},
				Command:     "/usr/bin/find",
			},
			wantErr: false,
		},
		{
			name:       "Complex cron - combinations",
			expression: "*/15 0 1-7,15,21-23/2 * 1-5 /usr/bin/find",
			want: &CronSchedule{
				Minutes:     []int{0, 15, 30, 45},
				Hours:       []int{0},
				DaysOfMonth: []int{1, 2, 3, 4, 5, 6, 7, 15, 21, 23},
				Months:      makeRange(1, 12),
				DaysOfWeek:  []int{1, 2, 3, 4, 5},
				Command:     "/usr/bin/find",
			},
			wantErr: false,
		},
		{
			name:       "Complex cron - multiple intervals",
			expression: "*/15 */6 1-15/2 */3 */2 /usr/bin/find",
			want: &CronSchedule{
				Minutes:     []int{0, 15, 30, 45},
				Hours:       []int{0, 6, 12, 18},
				DaysOfMonth: []int{1, 3, 5, 7, 9, 11, 13, 15},
				Months:      []int{1, 4, 7, 10},
				DaysOfWeek:  []int{0, 2, 4, 6},
				Command:     "/usr/bin/find",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.expression)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errorString != "" && err.Error() != tt.errorString {
					t.Errorf("Parse() error = %v, want %v", err.Error(), tt.errorString)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() mismatch:\ngot = %+v\nwant = %+v", got, tt.want)
			}
		})
	}
}

func makeRange(min, max int) []int {
	result := make([]int, max-min+1)
	for i := range result {
		result[i] = min + i
	}
	return result
}
