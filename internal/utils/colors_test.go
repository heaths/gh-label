package utils

import (
	"regexp"
	"testing"
)

func Test_RandomColor(t *testing.T) {
	t.Run("validate random colors", func(t *testing.T) {
		re := regexp.MustCompile("^[A-Z0-9]{6}$")
		for i := 0; i < 10; i++ {
			color := RandomColor()
			if !re.MatchString(color) {
				t.Errorf("RandomColor() = %s, want pattern: %s", color, re.String())
			}
		}
	})
}

func Test_ValidateColor(t *testing.T) {
	tests := []struct {
		name  string
		color string
		want  string
		wantE bool
	}{
		{
			name:  "valid color",
			color: "11bb33",
			want:  "11bb33",
		},
		{
			name:  "valid color with # prefix",
			color: "#11bb33",
			want:  "11bb33",
		},
		{
			name:  "color too short",
			color: "1b3",
			wantE: true,
		},
		{
			name:  "color too long",
			color: "11bb33dd",
			wantE: true,
		},
		{
			name:  "invalid hex digit",
			color: "aabbzz",
			wantE: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ValidateColor(tt.color); (err != nil) != tt.wantE {
				t.Errorf("ValidateColor() error = %v, wantE: %v", err, tt.wantE)
			} else if got != tt.want {
				t.Errorf("ValidateColor() color = %s, want: %s", got, tt.want)
			}
		})
	}
}
