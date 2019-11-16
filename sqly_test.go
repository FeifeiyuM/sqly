package sqlyt

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var opt = &Option{
	Dsn:             "test:mysql123@tcp(localhost:3306)/test_db?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local",
	DriverName:      "mysql",
	MaxIdleConns:    0,
	MaxOpenConns:    0,
	ConnMaxLifeTime: 0,
}

// user model
type Account struct {
	ID         int64      `sql:"id" json:"id"`
	Nickname   string     `sql:"nickname" json:"nickname"`
	Avatar     NullString `sql:"avatar" json:"avatar"`
	Email      string     `sql:"email" json:"email"`
	Mobile     string     `sql:"mobile" json:"mobile"`
	Role       int8       `sql:"role" json:"role"`
	Password   string     `sql:"password" json:"password"`
	CreateTime time.Time  `sql:"create_time" json:"create_time"`
}

func TestNew(t *testing.T) {
	_, err := New(opt)
	if err != nil {
		t.Error(err)
	}
}

func TestSqlY_Exec(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	//query := "CREATE DATABASE `test_db`;"
	//_, err = db.Exec(query)
	//if err != nil {
	//	t.Error(err)
	//}

	query := "CREATE TABLE `account` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`nickname` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`avatar` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'avatar url'," +
		"`mobile` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'mobile number'," +
		"`email` varchar(320) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'email'," +
		"`password` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'password'," +
		"`role` tinyint(4) DEFAULT '0' COMMENT 'role'," +
		"`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `mobile_index` (`mobile`)," +
		"KEY `email_index` (`email`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"
	_, err = db.Exec(query)
	if err != nil {
		t.Error(err)
	}
}

func TestSqlY_QueryOne(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}

	acc := new(Account)
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` FROM `account` " +
		"WHERE `id`=?;"
	err = db.QueryOne(acc, query, 6)
	if err != nil {
		t.Error(err)
	}
}

func TestSqlY_Query(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	acc := new(Account)
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` FROM `account` " +
		"WHERE `avatar`='';"
	rows, err := db.Query(acc, query, nil)
	if err != nil {
		t.Error(err)
	}
	res := map[string]interface{}{
		"acc": acc,
	}
	accStr, _ := json.Marshal(res)
	fmt.Printf("rows %s", accStr)
	_ = fmt.Sprint(rows)
}

func TestSqlY_Insert(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
		"VALUES (?, ?, ?, ?);"
	aff, err := db.Insert(query, "nick_test3", "18812311235", "test@foxmail.com", 1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(aff)
}

func TestSqlY_Update(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "UPDATE `account` SET `nickname`=? WHERE `mobile`=?;"
	aff, err := db.Update(query, "lucy", "18812311231")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(aff)
}

func TestSqlY_Delete(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "DELETE FROM `account` WHERE `mobile`=?;"
	aff, err := db.Delete(query, "18812311231")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(aff)
}

func TestSqlY_InsertCtx(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
		"VALUES (?, ?, ?, ?);"
	ctx := context.TODO()
	aff, err := db.InsertCtx(ctx, query, "nick_test2", "18812311232", "test2@foxmail.com", 1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(aff)
}

func TestSqlY_QueryCtx(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
		"FROM `account` WHERE `avatar` IS ?;"
	ctx := context.TODO()
	acc := new(Account)
	res, err := db.QueryCtx(ctx, acc, query, nil)
	if err != nil {
		t.Error(err)
	}
	resStr, _ := json.Marshal(res)
	fmt.Println(resStr)
}

func TestSqlY_QueryOneCtx(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role`, `create_time` " +
		"FROM `account` WHERE `mobile`=?;"
	ctx := context.TODO()
	acc := new(Account)
	err = db.QueryOneCtx(ctx, acc, query, "18756788776")
	if err != nil {
		t.Error(err)
	}
	resStr, _ := json.Marshal(acc)
	fmt.Println(resStr)
	acc2 := new(Account)
	err = json.Unmarshal(resStr, acc2)
	if err != nil {
		t.Error(err)
	}
}

func TestSqlY_InsertMany(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
		"VALUES (?, ?, ?, ?);"
	var vals = [][]interface{}{
		{"testq1", "18112342345", "testq1@foxmail.com", 1},
		{"testq2", "18112342346", "testq2@foxmail.com", 1},
	}
	aff, err := db.InsertMany(query, vals)
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintln(aff)
}

func TestSqlY_ExecManyCtx(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	ctx := context.TODO()
	var queries []string
	query, _ := QueryFmt("UPDATE `account` SET `nickname`=? WHERE `mobile`=?;", "nick_many", "18112342345")
	queries = append(queries, query)
	query, _ = QueryFmt("DELETE FROM `account` WHERE `mobile`=?;", "18112342346")
	queries = append(queries, query)
	err = db.ExecManyCtx(ctx, queries)
	if err != nil {
		t.Error(err)
	}
}
