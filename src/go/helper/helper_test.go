package helper

import "testing"

func TestStringSlice_Contains(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		s    StringSlice
		args args
		want bool
	}{
		{
			name: "success",
			s:    []string{"abc", "def", "ghi"},
			args: args{
				str: "def",
			},
			want: true,
		},
		{
			name: "error",
			s:    []string{"abc", "def", "ghi"},
			args: args{
				str: "xyz",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Contains(tt.args.str); got != tt.want {
				t.Errorf("StringSlice.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceContains(t *testing.T) {
	type args struct {
		s   []string
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{
				s:   []string{"abc", "def", "ghi"},
				str: "def",
			},
			want: true,
		},
		{
			name: "error",
			args: args{
				s:   []string{"abc", "def", "ghi"},
				str: "xyz",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceContains(tt.args.s, tt.args.str); got != tt.want {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSlice_Contains(t *testing.T) {
	type args struct {
		e int
	}
	tests := []struct {
		name string
		s    IntSlice
		args args
		want bool
	}{
		{
			name: "success",
			s:    []int{1, 2, 3},
			args: args{
				e: 2,
			},
			want: true,
		},
		{
			name: "error",
			s:    []int{1, 2, 3},
			args: args{
				e: 4,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Contains(tt.args.e); got != tt.want {
				t.Errorf("IntSlice.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntSliceContains(t *testing.T) {
	type args struct {
		s []int
		e int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "success",
			args: args{
				s: []int{1, 2, 3},
				e: 2,
			},
			want: true,
		},
		{
			name: "error",
			args: args{
				s: []int{1, 2, 3},
				e: 4,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntSliceContains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("IntSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
