package main

import "testing"

func TestSum(t *testing.T) {
	result := sum(1, 2)
	expected := 3

	if result != expected {
		t.Errorf("Sum returned %d, expected %d", result, expected)
	}
}
