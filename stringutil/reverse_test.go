package stringutil

import "testing"

type TestData struct{
	in, want string
}

func TestReverse(t *testing.T){
	cases:= []TestData{
		{"Hello world", "dlrow olleH"},
		{"this is a test", "tset a si siht"},
	}

	for _, c := range cases{
		got:=Reverse(c.in)
		if c.want != got{
			t.Errorf("Reverse(%q), get %q. Expect %q", c.in, got, c.want)
		}
	}
}
