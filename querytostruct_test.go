package querytostruct

import (
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type test struct {
		Str1            string `url:"str1"`
		StrBlock        string `url:"-"`
		StrNoTag        *string
		Strings         []string        `url:"str"`
		Bytes           []byte          `url:"bytes"`
		Int1            int             `url:"int1"`
		Uint16          uint16          `url:"uint16"`
		Float32         float32         `url:"float32"`
		Ints            []int           `url:"int"`
		NonExistentInts []int           `url:"nonint"`
		Bool1           bool            `url:"bool1"`
		Bools           []bool          `url:"bool"`
		NonExistent     []string        `url:"non"`
		BadNum          int             `url:"badnum"`
		BadNumSlice     []int           `url:"badnumslice"`
		OtherTag        string          `form:"otherval"`
		OmitEmpty       string          `form:"otherval,omitempty"`
		OtherTags       string          `url:"othertags" json:"othertags"`
		NotSupported    map[string]bool `url:"notsupported" json:"notsupported"`
	}

	q := url.Values{}
	q.Add("str1", "string1")
	q.Add("str", "str1")
	q.Add("str", "str2")
	q.Add("str", "str3")
	q.Add("bytes", "manybytes")
	q.Add("int1", "123")
	q.Add("uint16", "123")
	q.Add("float32", "123.456")
	q.Add("int", "456")
	q.Add("int", "789")
	q.Add("bool1", "true")
	q.Add("bool", "true")
	q.Add("bool", "false")
	q.Add("bool", "f")
	q.Add("bool", "t")
	q.Add("badnum", "abc")
	q.Add("badnumslice", "abc")
	q.Add("badnumslice", "def")
	q.Add("notsupported", "def")

	// Bad.
	var b map[string]interface{}
	if _, err := Unmarshal(q, &b, "url"); err == nil {
		t.Error("non-struct target passed")
	}

	// Good.
	var o test
	Unmarshal(q, &o, "url")
	exp := test{
		Str1:            "string1",
		Strings:         []string{"str1", "str2", "str3"},
		Bytes:           []byte("manybytes"),
		Int1:            123,
		Uint16:          123,
		Float32:         123.456,
		Ints:            []int{456, 789},
		NonExistentInts: nil,
		Bool1:           true,
		Bools:           []bool{true, false, false, true},
		BadNum:          0,
		BadNumSlice:     []int{0, 0},
	}
	if !reflect.DeepEqual(exp, o) {
		t.Error("scan structs don't match. expected != scanned")
		fmt.Println(exp)
		fmt.Println(o)
	}
}
