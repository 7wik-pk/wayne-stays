package config

import (
	"reflect"
	"testing"
)

func TestGetEnvSlice(t *testing.T) {
	defaults := []string{"http://localhost:5173"}

	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{"unset uses default", "", defaults},
		{"single value", "http://localhost:3000", []string{"http://localhost:3000"}},
		{"trims spaces", "http://a.test, http://b.test", []string{"http://a.test", "http://b.test"}},
		{"drops empty items", "http://a.test,,http://b.test,", []string{"http://a.test", "http://b.test"}},
		{"only separators uses default", " , ", defaults},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("TEST_SLICE", tt.value)
			if got := getEnvSlice("TEST_SLICE", defaults); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getEnvSlice() = %q, want %q", got, tt.want)
			}
		})
	}
}
