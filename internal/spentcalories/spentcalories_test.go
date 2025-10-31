package spentcalories

import (
	"testing"
)

func TestTrainingInfo_ValidWalking(t *testing.T) {
	result, err := TrainingInfo("1000,Ходьба,1h0m", 70.0, 1.75)
	if err != nil {
		t.Errorf("Expected no error for valid walking input, got: %v", err)
	}
	if result == "" {
		t.Error("Expected non-empty result for valid walking input")
	}
}

func TestTrainingInfo_ValidRunning(t *testing.T) {
	result, err := TrainingInfo("1000,Бег,1h0m", 70.0, 1.75)
	if err != nil {
		t.Errorf("Expected no error for valid running input, got: %v", err)
	}
	if result == "" {
		t.Error("Expected non-empty result for valid running input")
	}
}

func TestTrainingInfo_InvalidInput(t *testing.T) {
	_, err := TrainingInfo("invalid,data,1h", 70.0, 1.75)
	if err == nil {
		t.Error("Expected error for invalid input")
	}
}
