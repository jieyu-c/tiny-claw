package tools

import "testing"

func Test_fuzzyReplace(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		originalContent string
		oldText         string
		newText         string
		want            string
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := fuzzyReplace(tt.originalContent, tt.oldText, tt.newText)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("fuzzyReplace() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("fuzzyReplace() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("fuzzyReplace() = %v, want %v", got, tt.want)
			}
		})
	}
}
