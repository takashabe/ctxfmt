package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_containPartial(t *testing.T) {
	tests := []struct {
		ss     []string
		substr string
		want   bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{"a", "b", "c"}, "ab", true},
		{[]string{"Repository", "Service"}, "FooRepository", true},
	}
	for _, tt := range tests {
		got := containPartial(tt.ss, tt.substr)
		assert.Equal(t, tt.want, got, "case: %v", tt)
	}
}
