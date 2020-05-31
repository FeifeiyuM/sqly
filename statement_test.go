package sqly

import (
	"fmt"
	"testing"
	"time"
)

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

	query = "INSERT INTO `accounts`(`mobile`, `gender`, `age`, `balance`, `address`, `status`) VALUES " +
		"(?,?,?,?,?,?)"
	res, err = queryFormat(query, "18887655678", NullString{String: "male"}, NullInt32{Int32: 12, Valid: true},
		NullFloat64{}, NullString{}, NullBool{Bool: false, Valid: true})
	if err != nil {
		t.Error(err)
	}
	resCmp = "INSERT INTO `accounts`(`mobile`, `gender`, `age`, `balance`, `address`, `status`) VALUES (\"18887655678\",\"male\",12,NULL,NULL,0)"
	if res != resCmp {
		t.Error("error")
	}

	query = "INSERT INTO `accounts`(`create_time`, `add_time`) VALUES (?,?)"
	res, err = queryFormat(query, time.Now(), NullTime{Time: time.Now(), Valid: true})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)

	query = "SELECT * FROM `accounts` WHERE `add_time` IN ?"
	times := []time.Time{time.Now(), time.Now(), time.Now()}
	res, err = queryFormat(query, times)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(res)
}
