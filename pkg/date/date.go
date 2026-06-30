package date

import (
	"strings"
	"time"
)

type Date struct {
	time.Time
}

// New is ...
func New(date string, format string) (Date, error) {
	format = formatPattern(format)

	formattedDate, err := time.Parse(format, date)
	if err != nil {
		return Date{}, err
	}

	goment := Date{formattedDate}

	return goment, nil
}

func Now() Date {
	now := time.Now()
	return Date{now}
}

// Formats is ...
func (g Date) Formats(format string) string {
	format = formatPattern(format)

	return g.Format(format)
}

func formatPattern(format string) string {
	format = strings.ReplaceAll(format, "YYYY", "2006")
	format = strings.ReplaceAll(format, "YYY", "006")
	format = strings.ReplaceAll(format, "YY", "06")

	format = strings.ReplaceAll(format, "MMMM", "January")
	format = strings.ReplaceAll(format, "MMM", "Jan")
	format = strings.ReplaceAll(format, "MM", "01")
	format = strings.ReplaceAll(format, "M", "1")

	format = strings.ReplaceAll(format, "DD", "02")
	format = strings.ReplaceAll(format, "D", "2")

	format = strings.ReplaceAll(format, "dddd", "Monday")
	format = strings.ReplaceAll(format, "ddd", "Mon")

	format = strings.ReplaceAll(format, "HH", "15")
	format = strings.ReplaceAll(format, "hh", "03")
	format = strings.ReplaceAll(format, "h", "3")

	format = strings.ReplaceAll(format, "A", "PM")

	format = strings.ReplaceAll(format, "mm", "04")
	format = strings.ReplaceAll(format, "m", "4")

	format = strings.ReplaceAll(format, "ss", "05")
	format = strings.ReplaceAll(format, "s", "5")

	format = strings.ReplaceAll(format, "SSSSSS", "999999")

	format = strings.ReplaceAll(format, "zz", "MST")
	format = strings.ReplaceAll(format, "z", "MST")

	format = strings.ReplaceAll(format, "ZZ", "Z0700")

	return format
}
