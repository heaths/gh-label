package utils

import (
	"fmt"
	"testing"
)

func TestAreEqualStrings(t *testing.T) {
	tests := []struct {
		a    []string
		b    []string
		want bool
	}{
		{
			a:    nil,
			b:    nil,
			want: true,
		},
		{
			a:    []string{},
			b:    []string{},
			want: true,
		},
		{
			a:    []string{"x"},
			b:    []string{"x"},
			want: true,
		},
		{
			a:    []string{"x", "y"},
			b:    []string{"x", "y"},
			want: true,
		},
		{
			a:    []string{"x"},
			b:    []string{"y"},
			want: false,
		},
		{
			a:    []string{"x"},
			b:    []string{"x", "y"},
			want: false,
		},
		{
			a:    []string{"x"},
			b:    nil,
			want: false,
		},
		{
			a:    []string{"x"},
			b:    []string{},
			want: false,
		},
		{
			a:    []string{"y"},
			b:    []string{"x"},
			want: false,
		},
		{
			a:    []string{"x", "y"},
			b:    []string{"x"},
			want: false,
		},
		{
			a:    nil,
			b:    []string{"x"},
			want: false,
		},
		{
			a:    []string{},
			b:    []string{"x"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("a = %v, b = %v", tt.a, tt.b), func(t *testing.T) {
			if got := AreEqualStrings(tt.a, tt.b); got != tt.want {
				t.Errorf("AreEqualStrings() = %v, want %v", got, tt.want)
			}
		})
	}
}
