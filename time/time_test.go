package time

import (
	"strings"
	"testing"

	"github.com/gxxgle/go-utils/json"
)

func TestJSON(t *testing.T) {
	before := New(2008, M8, 8, 20, 8, 8, 8*1000*1000)
	str := json.MustMarshalToString(before)
	after := Time{}

	if err := json.UnmarshalFromString(str, &after); err != nil {
		t.Fatal("time.Time UnmarshalJSON failed, err:", err)
	}

	if !strings.HasPrefix(str, "\"2008-08-08T20:08:08") {
		t.Fatal("time.Time MarshalJSON not match, str:", str)
	}

	if before != after {
		t.Fatal("time.Time UnmarshalJSON not match")
	}
}
