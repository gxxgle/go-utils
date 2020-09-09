package json

import (
	"testing"
)

func TestFormat(t *testing.T) {
	type check struct {
		raw  string
		want string
	}

	checks := []check{
		{raw: `  {  "a": "b"    } `, want: `{"a":"b"}`},
		{raw: "[\n1,\n2,\n3,\n4,\n5\n]", want: `[1,2,3,4,5]`},
	}

	for _, c := range checks {
		if format := FormatFromString(c.raw); format != c.want {
			t.Fatalf("TestFormat want: `%v`, but: `%v`", c.want, format)
		}
	}
}
