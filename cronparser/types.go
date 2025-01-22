package cronparser

type CronSchedule struct {
	Minutes     []int
	Hours       []int
	DaysOfMonth []int
	Months      []int
	DaysOfWeek  []int
	Command     string
}

type TimeField int

const (
	MinuteField TimeField = iota
	HourField
	DayOfMonthField
	MonthField
	DayOfWeekField
)

type ExpressionParser interface {
	ParseField(field string, timeField TimeField) ([]int, error)
}
