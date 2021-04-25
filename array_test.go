package sqly

import "testing"

func TestParseArray(t *testing.T) {
	//str := []byte("{10001, 10002, 10003, 10004}")
	str := []byte("{{\"meeting\",\"lunch\",\"lunch2\"},{\"training\",\"presentation\",\"fff\"}}")
	dims, elems, err := parseArray(str, []byte{','})
	if err != nil {
		t.Log(err)
	}
	t.Log(dims)
	t.Log(elems)
}
