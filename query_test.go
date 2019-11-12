package sqlyt

import "testing"

func TestQueryFmt(t *testing.T) {
	query := "select * from `accounts` WHERE `id`=? AND `status`=?;"
	res, err := QueryFmt(query, 9, 1)
	if err != nil {
		t.Error(err)
	}
	resCmp := "select * from `accounts` WHERE `id`=9 AND `status`=1;"
	if res != resCmp {
		t.Error("error")
	}

	query = "SELECT * FROM `accounts` WHERE `mobile`=? AND `role` IN ?;"
	res, err = QueryFmt(query, "18712342345", []int64{0, 1, 2})
	resCmp = "SELECT * FROM `accounts` WHERE `mobile`=\"18712342345\" AND `role` IN (0,1,2);"
	if err != nil {
		t.Error(err)
	}
	if res != resCmp {
		t.Error("error")
	}
}
