package main

import (
	"Deliveroo/cronparser"
	"fmt"
	"os"
	"strings"
)

func formatField(name string, values []int) string {
	paddedName := fmt.Sprintf("%-14s", name)
	valueStrs := make([]string, len(values))
	for i, v := range values {
		valueStrs[i] = fmt.Sprintf("%d", v)
	}
	return paddedName + strings.Join(valueStrs, " ")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go \"<cron expression>\"")
		os.Exit(1)
	}

	expression, err := cronparser.Parse(os.Args[1])

	if err != nil {
		fmt.Printf("Error parsing cron expression: %v\n", err)
		os.Exit(1)

	}

	fmt.Println(formatField("minute", expression.Minutes))
	fmt.Println(formatField("hour", expression.Hours))
	fmt.Println(formatField("day of month", expression.DaysOfMonth))
	fmt.Println(formatField("month", expression.Months))
	fmt.Println(formatField("day of week", expression.DaysOfWeek))
	fmt.Println(formatField("command", []int{}) + expression.Command)
}
