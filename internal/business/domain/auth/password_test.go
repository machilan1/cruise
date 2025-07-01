package auth

import (
	"testing"
)

func TestUser_ParsePassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		want     Password
		wantErr  bool
	}{
		// valid
		{
			name:     "8 characters - mixed",
			password: "12#$abAB",
			want:     Password{value: "12#$abAB"},
		},
		{
			name:     "20 characters - mixed",
			password: "12345abcde#$%^&ABCDE",
			want:     Password{value: "12345abcde#$%^&ABCDE"},
		},
		// length
		{
			name:     "7 characters",
			password: "12#abAB",
			wantErr:  true,
		},
		{
			name:     "65 characters",
			password: "123456abcdefgABCDEFG%123456abcdefgABCDEFG%123456abcdefgABCDEFG%12",
			wantErr:  true,
		},
		// spaces
		{
			name:     "space",
			password: " ",
			wantErr:  true,
		},
		{
			name:     "between spaces",
			password: "     16fe$7890    ",
			wantErr:  true,
		},
		{
			name:     "spaces between",
			password: "123 asd $f4",
			wantErr:  true,
		},
		// omit character type
		{
			name:     "no numeric",
			password: "abcde#$%^&ABCDE",
			wantErr:  true,
		},
		{
			name:     "no alphabetic",
			password: "12345#$%^&",
			wantErr:  true,
		},
		// additional edge cases
		{
			name:     "only special characters",
			password: "@$!%*#?&_-",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "only letters",
			password: "abcdefgh",
			wantErr:  true,
		},
		{
			name:     "only numbers",
			password: "12345678",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParsePassword(tt.password)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("failed to parse password: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("parser failed to reject invalid input: %q", tt.password)
				return
			}
			if tt.want != got {
				t.Errorf("expected %q ,got: %q", tt.password, got)
			}
		})
	}
}
