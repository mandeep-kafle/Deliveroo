package cronparser

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type WildcardParser struct {
	*TimeConstraints
}

func NewWildCardParser(tc *TimeConstraints) *WildcardParser {
	return &WildcardParser{
		TimeConstraints: tc,
	}
}

func (p *WildcardParser) ParseField(field string, timeField TimeField) ([]int, error) {
	min, max := p.GetBoundaries(timeField)
	result := make([]int, max-min+1)
	for i := range result {
		result[i] = min + i
	}
	return result, nil
}

type ValueListParser struct {
	*TimeConstraints
	intervalParser    ExpressionParser
	rangeParser       ExpressionParser
	singleValueParser ExpressionParser
}

func NewValueListParser(tc *TimeConstraints) *ValueListParser {
	return &ValueListParser{
		TimeConstraints:   tc,
		intervalParser:    &IntervalParser{tc},
		rangeParser:       &RangeParser{tc},
		singleValueParser: &SingleValueParser{tc},
	}
}

func (p *ValueListParser) ParseField(field string, timeField TimeField) ([]int, error) {
	values := strings.Split(field, ",")
	result := make([]int, 0)
	uniqueValues := make(map[int]bool) //need uniqueness

	for _, value := range values {
		var nums []int
		var err error

		if strings.Contains(value, "/") {
			nums, err = p.intervalParser.ParseField(value, timeField)
		} else if strings.Contains(value, "-") {
			nums, err = p.rangeParser.ParseField(value, timeField)
		} else {
			nums, err = p.singleValueParser.ParseField(value, timeField)
		}

		if err != nil {
			return nil, fmt.Errorf("invalid list item '%s': %w", value, err)
		}

		for _, num := range nums {
			if !uniqueValues[num] {
				result = append(result, num)
				uniqueValues[num] = true
			}
		}
	}

	sort.Ints(result)
	return result, nil
}

type IntervalParser struct {
	*TimeConstraints
}

func NewIntervalParser(tc *TimeConstraints) *IntervalParser {
	return &IntervalParser{
		TimeConstraints: tc,
	}
}

func (p *IntervalParser) ParseField(field string, timeField TimeField) ([]int, error) {
	parts := strings.Split(field, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid interval format: %s", field)
	}

	min, max := p.GetBoundaries(timeField)
	var start, end int

	if parts[0] == "*" {
		start, end = min, max
	} else {
		rangeValues, err := (&RangeParser{p.TimeConstraints}).ParseField(parts[0], timeField)
		if err != nil {
			return nil, err
		}
		start, end = rangeValues[0], rangeValues[len(rangeValues)-1]
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid interval value: %s", parts[1])
	}
	if interval <= 0 {
		return nil, fmt.Errorf("interval must be positive: %d", interval)
	}

	result := make([]int, 0, (end-start)/interval+1)
	for i := start; i <= end; i += interval {
		result = append(result, i)
	}
	return result, nil
}

type RangeParser struct {
	*TimeConstraints
}

func NewRangeParser(tc *TimeConstraints) *RangeParser {
	return &RangeParser{
		TimeConstraints: tc,
	}
}

func (p *RangeParser) ParseField(field string, timeField TimeField) ([]int, error) {
	bounds := strings.Split(field, "-")
	if len(bounds) != 2 {
		return nil, fmt.Errorf("invalid range format: %s", field)
	}

	start, err := strconv.Atoi(bounds[0])
	if err != nil {
		return nil, fmt.Errorf("invalid range start: %s", bounds[0])
	}

	end, err := strconv.Atoi(bounds[1])
	if err != nil {
		return nil, fmt.Errorf("invalid range end: %s", bounds[1])
	}

	if !p.IsValueInRange(start, timeField) || !p.IsValueInRange(end, timeField) {
		min, max := p.GetBoundaries(timeField)
		return nil, fmt.Errorf("range values must be between %d and %d", min, max)
	}

	if start > end {
		return nil, fmt.Errorf("range start cannot be greater than end: %d-%d", start, end)
	}

	result := make([]int, end-start+1)
	for i := range result {
		result[i] = start + i
	}
	return result, nil
}

type SingleValueParser struct {
	*TimeConstraints
}

func NewSingleValueParser(tc *TimeConstraints) *SingleValueParser {
	return &SingleValueParser{
		TimeConstraints: tc,
	}
}

func (p *SingleValueParser) ParseField(field string, timeField TimeField) ([]int, error) {
	value, err := strconv.Atoi(field)
	if err != nil {
		return nil, fmt.Errorf("invalid numeric value: %s", field)
	}

	if !p.IsValueInRange(value, timeField) {
		min, max := p.GetBoundaries(timeField)
		return nil, fmt.Errorf("value must be between %d and %d", min, max)
	}

	return []int{value}, nil
}
