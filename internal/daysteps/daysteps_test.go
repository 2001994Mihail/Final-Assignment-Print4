package daysteps

import (
	"testing"
)

func TestDayActionInfo_ValidInput(t *testing.T) {
	result := DayActionInfo("1000,1h0m", 70.0, 1.75)
	if result == "" {
		t.Error("Expected non-empty result for valid input")
	}
}

func TestDayActionInfo_InvalidInput(t *testing.T) {
	result := DayActionInfo("invalid,data", 70.0, 1.75)
	if result != "" {
		t.Error("Expected empty result for invalid input")
	}
}
