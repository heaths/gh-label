package utils

import "testing"

func Test_ColorE(t *testing.T) {
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
			if got, err := ColorE(tt.color); (err != nil) != tt.wantE {
				t.Errorf("ColorE() error = %v, wantE: %v", err, tt.wantE)
			} else if got != tt.want {
				t.Errorf("ColorE() color = %s, want: %s", got, tt.want)
			}
		})
	}
}
