package cronparser

import (
	"errors"
	"fmt"
	"strings"
)

type CronExpressionParser struct {
	timeConstraints *TimeConstraints
	fieldParsers    map[string]ExpressionParser
}


func (p *CronExpressionParser) parseTimeField(field string, timeField TimeField) ([]int, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return nil, errors.New("empty field not allowed")
	}

	if field == "*" {
		return p.fieldParsers["*"].ParseField(field, timeField)
	}

	if strings.Contains(field, ",") {
		return p.fieldParsers[","].ParseField(field, timeField)
	}

	if strings.Contains(field, "/") {
		return p.fieldParsers["/"].ParseField(field, timeField)
	}

	if strings.Contains(field, "-") {
		return p.fieldParsers["-"].ParseField(field, timeField)
	}

	return p.fieldParsers[""].ParseField(field, timeField)
}
func Parse(expression string) (*CronSchedule, error) {
	constraints := NewTimeConstraints()

	singleValueParser := NewSingleValueParser(constraints)
	rangeParser := NewRangeParser(constraints)
	intervalParser := NewIntervalParser(constraints)
	wildcardParser := NewWildCardParser(constraints)
	valueListParser := NewValueListParser(constraints)

	p := &CronExpressionParser{
		timeConstraints: constraints,
		fieldParsers: map[string]ExpressionParser{
			"*": wildcardParser,
			",": valueListParser,
			"/": intervalParser,
			"-": rangeParser,
			"":  singleValueParser,
		},
	}

	fields := strings.Fields(expression)
	if len(fields) < 6 {
		return nil, errors.New("invalid cron expression: requires at least 6 fields")
	}

	schedule := &CronSchedule{
		Command: strings.Join(fields[5:], " "),
	}

	scheduleFields := []*[]int{
		&schedule.Minutes,
		&schedule.Hours,
		&schedule.DaysOfMonth,
		&schedule.Months,
		&schedule.DaysOfWeek,
	}

	for i := 0; i < 5; i++ {
		values, err := p.parseTimeField(fields[i], TimeField(i))
		if err != nil {
			return nil, fmt.Errorf("error in field %d (%s): %w", i+1, fields[i], err)
		}
		*scheduleFields[i] = values
	}

	return schedule, nil
}
