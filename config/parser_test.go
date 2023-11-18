package config

import (
	"os"
	"testing"
)

func TestGetEnvWithFallback(t *testing.T) {
	// Set up test cases
	tests := []struct {
		key      string
		fallback string
		envValue string
		expected string
	}{
		// Test when the environment variable is set
		{"EXISTING_VAR", "fallback_value", "env_value", "env_value"},

		// Test when the environment variable is not set
		{"NON_EXISTENT_VAR", "fallback_value", "", "fallback_value"},
	}

	// Run the test cases
	for _, test := range tests {
		// Set the environment variable
		if test.envValue != "" {
			os.Setenv(test.key, test.envValue)
		} else {
			os.Unsetenv(test.key)
		}

		// Call the function
		result := GetEnvWithFallback(test.key, test.fallback)

		// Check if the result matches the expected value
		if result != test.expected {
			t.Errorf("GetEnvWithFallback(%s, %s) returned %s, expected %s", test.key, test.fallback, result, test.expected)
		}
	}
}
