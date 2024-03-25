package utils

import "testing"

func TestParseSize(t *testing.T) {
	type test struct {
		input string
		want  int64
	}
	tests := []test{
		{"5B", 5},
		{"100KB", 102400},
		{"1MB", 1048576},
		{"2MB", 2097152},
		{"1GB", 1073741824},
		{"fdd", 104857600},
	}
	for _, i := range tests {
		if result := ParseSize(i.input); result != i.want {
			t.Errorf("%s expected %d, but %d got", i.input, i.want, result)
		}
	}

}
