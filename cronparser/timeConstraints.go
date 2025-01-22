package cronparser

// base
type TimeConstraints struct {
	fieldBoundaries map[TimeField][2]int
}

func NewTimeConstraints() *TimeConstraints {
	return &TimeConstraints{
		fieldBoundaries: map[TimeField][2]int{
			MinuteField:     {0, 59},
			HourField:       {0, 23},
			DayOfMonthField: {1, 31},
			MonthField:      {1, 12},
			DayOfWeekField:  {0, 6},
		},
	}
}

func (tc *TimeConstraints) GetBoundaries(timeField TimeField) (min, max int) {
	boundaries := tc.fieldBoundaries[timeField]
	return boundaries[0], boundaries[1]
}

func (tc *TimeConstraints) IsValueInRange(value int, timeField TimeField) bool {
	min, max := tc.GetBoundaries(timeField)
	return value >= min && value <= max
}
