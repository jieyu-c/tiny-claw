package engine_test

import (
	"testing"
	"tiny-claw/internal/engine"
)

func TestClazzToSchema(t *testing.T) {
	tests := []struct {
		name string // description of this test case
	}{
		// TODO: Add test cases.
		{
			name: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine.ClazzToSchema()
		})
	}
}
