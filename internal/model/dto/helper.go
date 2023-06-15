package dto

import (
	"fmt"
	"log"
)

// pembantu untuk mengonversi string ke int
func ConvertToInt(str string) int {
	var value int
	_, err := fmt.Sscanf(str, "%d", &value)
	if err != nil {
		log.Fatal("Invalid ID provided")
	}

	return value
}

// pembantu untuk memvalidasi input yang tidak kosong.
func ValidateNonEmptyInput(input string) error {
	if input == "" {
		return fmt.Errorf("input must not be empty")
	}
	return nil
}
