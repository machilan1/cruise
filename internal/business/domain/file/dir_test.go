package file

import "testing"

func TestDir_NewDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		want    Dir
		wantErr bool
	}{
		{
			name:  "empty",
			value: "",
			want:  Dir{value: ""},
		},
		{
			name:  "root",
			value: "/",
			want:  Dir{value: ""},
		},
		{
			name:  "dot",
			value: ".",
			want:  Dir{value: ""},
		},
		{
			name:  "1 level",
			value: "dir",
			want:  Dir{value: "dir"},
		},
		{
			name:  "2 levels",
			value: "dir/subdir",
			want:  Dir{value: "dir/subdir"},
		},
		{
			name:  "with slash at the beginning",
			value: "/dir/subdir",
			want:  Dir{value: "dir/subdir"},
		},
		{
			name:  "with slash at the end",
			value: "dir/subdir/",
			want:  Dir{value: "dir/subdir"},
		},
		{
			name:  "with slash at the beginning and end",
			value: "/dir/subdir/",
			want:  Dir{value: "dir/subdir"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewDir(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("NewDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDir_ParseDir(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		want    Dir
		wantErr bool
	}{
		{
			name: "empty",
			path: "",
			want: Dir{value: ""},
		},
		{
			name: "root",
			path: "/",
			want: Dir{value: ""},
		},
		{
			name: "dot",
			path: ".",
			want: Dir{value: ""},
		},
		{
			name: "1 level",
			path: "dir/file.jpg",
			want: Dir{value: "dir"},
		},
		{
			name: "2 levels",
			path: "dir/subdir/file.jpg",
			want: Dir{value: "dir/subdir"},
		},
		{
			name: "with slash at the beginning",
			path: "/dir/subdir/file.jpg",
			want: Dir{value: "dir/subdir"},
		},
		{
			name: "without file",
			path: "dir/subdir/",
			want: Dir{value: "dir/subdir"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseDirFromPath(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDirFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("ParseDirFromPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
