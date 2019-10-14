package utils

import (
	"testing"
)

func TestThatItCanRequireFields(t *testing.T) {
	parser := NewRequiredFlagSet([]string{"name"})

	var _ *string = parser.String("name", "", "Person's name")
	var _ *int = parser.Int("age", 100, "Person's age")

	err := parser.parseArguments([]string{
		"-age=18",
	})

	if err == nil {
		t.Error("Should have returned an error")
	} else {
		if err.Error() != "name is a required flag" {
			t.Errorf("Should have returned a required error, got %s", err.Error())
		}
	}

	err = parser.parseArguments([]string{
		"-age=18",
		"-name=Bob",
	})

	if err != nil {
		t.Errorf("Should not have received error, instead got %s", err.Error())
	}
}