package webhook

import (
	"testing"
)

func TestExtractJSONPath(t *testing.T) {
	tests := []struct {
		name     string
		body     []byte
		path     string
		expected string
	}{
		{
			name:     "simple string",
			body:     []byte(`{"record_id": "abc123"}`),
			path:     "record_id",
			expected: "abc123",
		},
		{
			name:     "nested path",
			body:     []byte(`{"data": {"record_id": "nested-456"}}`),
			path:     "data.record_id",
			expected: "nested-456",
		},
		{
			name:     "deeply nested",
			body:     []byte(`{"a": {"b": {"c": "deep-value"}}}`),
			path:     "a.b.c",
			expected: "deep-value",
		},
		{
			name:     "missing key",
			body:     []byte(`{"data": {"other": "value"}}`),
			path:     "data.record_id",
			expected: "",
		},
		{
			name:     "invalid json",
			body:     []byte(`not json`),
			path:     "data.record_id",
			expected: "",
		},
		{
			name:     "empty body",
			body:     []byte(``),
			path:     "data.record_id",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractJSONPath(tt.body, tt.path)
			if result != tt.expected {
				t.Errorf("extractJSONPath(%s, %s) = %q, want %q", string(tt.body), tt.path, result, tt.expected)
			}
		})
	}
}
