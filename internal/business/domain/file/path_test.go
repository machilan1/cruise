package file

import "testing"

func TestPath_String(t *testing.T) {
	tests := []struct {
		name string
		p    Path
		want string
	}{
		{
			name: "simple",
			p:    Path{Dir{value: "dir"}, "base", "ext"},
			want: "dir/base.ext",
		},
		{
			name: "no dir",
			p:    Path{Dir{}, "base", "ext"},
			want: "base.ext",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.p.String(); got != tt.want {
				t.Errorf("Path.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
