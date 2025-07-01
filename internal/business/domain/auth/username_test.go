package auth

import (
	"testing"
)

func TestUser_ParseUserName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		username string
		want     Username
		wantErr  bool
	}{
		// valid
		{
			name:     "4 character - alphabetic",
			username: "wxyz",
			want:     Username{value: "wxyz"},
		},
		{
			name:     "4 character - numeric",
			username: "0123",
			want:     Username{value: "0123"},
		},
		{
			name:     "30 characters - mixed",
			username: "1234567890abcdefghijABCDEFGHIJ",
			want:     Username{"1234567890abcdefghijABCDEFGHIJ"},
		},
		// length
		{
			name:     "0 character",
			username: "",
			wantErr:  true,
		},
		{
			name:     "3 characters",
			username: "abc",
			wantErr:  true,
		},
		{
			name:     "31 characters - mixed",
			username: "1234567890abcdefghijABCDEFGHIJx",
			wantErr:  true,
		},
		// spaces
		{
			name:     "space",
			username: " ",
			wantErr:  true,
		},
		{
			name:     "between spaces",
			username: "     1234567890    ",
			wantErr:  true,
		},
		{
			name:     "space between",
			username: "123 asd 4f4",
			wantErr:  true,
		},
		// special characters
		{
			name:     "email",
			username: "johndoe@example.com",
			wantErr:  true,
		},
		{
			name:     "special characters",
			username: "S@#D$$D@^&",
			wantErr:  true,
		},
		{
			name:     "only 1 special character",
			username: "1234&abcd4",
			wantErr:  true,
		},
		{
			name:     "non ASCII",
			username: "用戶名",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseUsername(tt.username)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("failed to parse username: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("parser failed to reject given input: %q", tt.username)
				return
			}
			if got != tt.want {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
