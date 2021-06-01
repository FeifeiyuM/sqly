package sqly

import (
	"fmt"
	"testing"
	"time"
)

func TestQueryFmt(t *testing.T) {
	query := "select * from `accounts` WHERE `id`=? AND `status`=?;"
	res, err := QueryFmtMysql(query, 9, 1)
	if err != nil {
		t.Error(err)
	}
	resCmp := "select * from `accounts` WHERE `id`=9 AND `status`=1;"
	if res != resCmp {
		t.Error("error")
	}

	query = "SELECT * FROM `accounts` WHERE `mobile`=? AND `role` IN ?;"
	res, err = QueryFmtMysql(query, "18712342345", []int64{0, 1, 2})
	resCmp = "SELECT * FROM `accounts` WHERE `mobile`='18712342345' AND `role` IN (0,1,2);"
	if err != nil {
		t.Error(err)
	}
	if res != resCmp {
		t.Error("error")
	}

	query = "INSERT INTO `accounts`(`mobile`, `gender`, `age`, `balance`, `address`, `status`) VALUES " +
		"(?,?,?,?,?,?)"
	res, err = QueryFmtMysql(query, "18887655678", NullString{String: "male"}, NullInt32{Int32: 12, Valid: true},
		NullFloat64{}, NullString{}, NullBool{Bool: false, Valid: true})
	if err != nil {
		t.Error(err)
	}
	resCmp = "INSERT INTO `accounts`(`mobile`, `gender`, `age`, `balance`, `address`, `status`) VALUES ('18887655678','male',12,NULL,NULL,0)"
	if res != resCmp {
		t.Error("error")
	}

	query = "INSERT INTO `accounts`(`create_time`, `add_time`) VALUES (?,?)"
	res, err = QueryFmtMysql(query, time.Now(), NullTime{Time: time.Now(), Valid: true})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)

	query = "SELECT * FROM `accounts` WHERE `add_time` IN ?"
	times := []time.Time{time.Now(), time.Now(), time.Now()}
	res, err = QueryFmtMysql(query, times)
	if err != nil {
		t.Error(err)
	}
	fmt.Print(res)

	// array
	query = "INSERT INTO `accounts` (`sub_ids`) VALUES (?)"
	ids := []int64{1, 2, 12345, 34}
	res, err = QueryFmtMysql(query, Array(ids))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
	query = "UPDATE `accounts` SET `email`=? WHERE `id`=?"
	emails := []string{"t1@qq.com", "t2@qq.com", "t3@qq.com"}
	res, err = QueryFmtMysql(query, Array(emails), 11)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}

func TestQueryFmt_pg(t *testing.T) {
	query := "UPDATE `accounts` SET `email`=? WHERE `id`=?"
	emails := []string{"t1@qq.com", "t2@qq.com", "t3@qq.com"}
	res, err := QueryFmtPostgresql(query, Array(emails), 11)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}

func TestQueryFmtPostgresql(t *testing.T) {
	errMsg := "rpc error: code = InvalidArgument desc = 参数错误: error=Key: 'CommonOrderRequest.OrderNumber' Error:Field validation for 'OrderNumber' failed on the 'required' tag"
	query := "UPDATE schedules SET err_msg=? WHERE id=?"
	res, err := QueryFmtPostgresql(query, errMsg, 1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
