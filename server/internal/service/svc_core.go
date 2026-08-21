package service


import (
	"strings"
	"time"
)


func parse1CDate(dateStr string) *time.Time {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" || dateStr == "00.00.0000" || dateStr == "01-01-0001" {
		return nil
	}
	t, err := time.Parse("02-01-2006", dateStr)
	if err != nil {
		return nil
	}
	return &t
}

func stringPtr(str string) *string {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}
	return &str
}

func logSkip(name, reason string) map[string]string {
	return map[string]string{"item": name, "reason": reason}
}

