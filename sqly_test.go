package sqly

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var opt = &Option{
	Dsn:             "root:root@tcp(127.0.0.1:3306)/test_db?multiStatements=true&charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local",
	DriverName:      "mysql",
	MaxIdleConns:    0,
	MaxOpenConns:    0,
	ConnMaxLifeTime: 0,
}

// user model
type Account struct {
	ID         int64       `sql:"id" json:"id"`
	Nickname   string      `sql:"nickname" json:"nickname"`
	Avatar     NullString  `sql:"avatar" json:"avatar"`
	Email      string      `sql:"email" json:"email"`
	Mobile     string      `sql:"mobile" json:"mobile"`
	Role       NullInt32   `sql:"role" json:"role"`
	Password   string      `sql:"password" json:"password"`
	IsValid    NullBool    `sql:"is_valid" json:"is_valid"`
	Stature    NullFloat64 `sql:"stature" json:"stature"`
	CreateTime time.Time   `sql:"create_time" json:"create_time"`
	AddTime    NullTime    `sql:"add_time" json:"add_time"`
	Birthday   NullTime    `sql:"birthday" json:"birthday"`
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

	query := "DROP TABLE IF EXISTS `account`;" +
		"CREATE TABLE `account` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`nickname` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL," +
		"`avatar` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'avatar url'," +
		"`mobile` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'mobile number'," +
		"`email` varchar(320) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'email'," +
		"`password` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'password'," +
		"`role` tinyint(4) DEFAULT '0' COMMENT 'role'," +
		"`is_valid` tinyint(4) DEFAULT NULL COMMENT 'is_valid'," +
		"`stature` float(5,2) DEFAULT NULL COMMENT 'stature'," +
		"`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"`add_time` datetime DEFAULT NULL, " +
		"`birthday` date DEFAULT NULL, " +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `mobile_index` (`mobile`)," +
		"KEY `email_index` (`email`)" +
		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"
	_, err = db.Exec(query)
	if err != nil {
		t.Error(err)
	}
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

func TestSqlY_QueryOne(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}

	acc := new(Account)
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` FROM `account` " +
		"WHERE `id`=?;"
	err = db.Get(acc, query, 1)
	if err != nil {
		t.Error(err)
	}
}

func TestSqlY_Query(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	var accs []*Account
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role`, `create_time` FROM `account`;"

	err = db.Query(&accs, query, nil)
	if err != nil {
		t.Error(err)
	}
	accStr, _ := json.Marshal(accs)
	fmt.Printf("rows %s", accStr)
}

func TestSqlY_Query_All(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	var accs []*Account
	query := "SELECT * FROM `account` limit 1;"

	err = db.Query(&accs, query, nil)
	if err != nil {
		t.Error(err)
	}
	accStr, _ := json.Marshal(accs)
	fmt.Printf("rows %s", accStr)
}

func TestSqlY_Get_All(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	acc := &Account{}
	query := "SELECT * FROM `account` WHERE `id`=2"

	err = db.Get(acc, query, nil)
	if err != nil {
		t.Error(err)
	}
	accStr, _ := json.Marshal(acc)
	fmt.Printf("rows %s", accStr)
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

func TestSqlY_QueryCtx(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
		"FROM `account` WHERE `avatar` IS ?;"
	ctx := context.TODO()
	var acc []Account
	err = db.QueryCtx(ctx, &acc, query, nil)
	if err != nil {
		t.Error(err)
	}
	resStr, _ := json.Marshal(acc)
	fmt.Println(resStr)
}

func TestSqlY_GetCtx(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role`, `create_time` " +
		"FROM `account` WHERE `mobile`=?;"
	ctx := context.TODO()
	acc := new(Account)
	err = db.GetCtx(ctx, acc, query, "18812311232")
	if err != nil {
		t.Error(err)
	}
	resStr, _ := json.Marshal(acc)
	fmt.Println(string(resStr))
	acc2 := new(Account)
	err = json.Unmarshal(resStr, acc2)
	if err != nil {
		t.Error(err)
	}
}

func TestSqlY_GetCtx_Empty(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role`, `create_time` " +
		"FROM `account` WHERE `mobile`=?;"
	ctx := context.TODO()
	acc := new(Account)
	err = db.GetCtx(ctx, acc, query, "18812311239")
	if err != ErrEmpty {
		t.Error("expect error empty")
	}
}

func TestSqlY_GetCtx_Multi(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role`, `create_time` " +
		"FROM `account`;"
	ctx := context.TODO()
	acc := new(Account)
	err = db.GetCtx(ctx, acc, query)
	if err != ErrMultiRes {
		t.Error("expect multi results error")
	}
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

func TestSqlY_NewTrans(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
		return
	}
	ts, err := db.NewTrans()
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		_ = ts.Rollback()
	}()

	ctx := context.TODO()

	acc := new(Account)
	query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role`, `create_time` " +
		"FROM `account` WHERE `mobile`=?;"
	err = ts.GetCtx(ctx, acc, query, "18812311235")
	if err != nil {
		t.Error(err)
		return
	}
	query = "UPDATE `account` SET `nickname`=? WHERE `id`=?;"
	name := `<ok class="12">`
	aff, err := ts.UpdateCtx(ctx, query, name, acc.ID)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(aff)

	err = ts.Commit()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestBoolean_Scan(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
		return
	}

	type Acc struct {
		ID         int64       `sql:"id" json:"id"`
		Nickname   string      `sql:"nickname" json:"nickname"`
		Avatar     NullString  `sql:"avatar" json:"avatar"`
		Email      string      `sql:"email" json:"email"`
		Mobile     string      `sql:"mobile" json:"mobile"`
		Role       NullInt32   `sql:"role" json:"role"`
		Password   string      `sql:"password" json:"password"`
		IsValid    NullBool    `sql:"is_valid" json:"is_valid"`
		Stature    NullFloat64 `sql:"stature" json:"stature"`
		CreateTime time.Time   `sql:"create_time" json:"create_time"`
	}
	var accs []*Acc
	query := "SELECT `id`, `nickname`,  `mobile`, `password`, `role`, `create_time`, `is_valid` FROM `account`;"
	err = db.Query(&accs, query)
	if err != nil {
		t.Error(err)
	}
	resStr, _ := json.Marshal(accs)
	fmt.Println(string(resStr))
	var accs2 []*Acc
	err = json.Unmarshal(resStr, &accs2)
	if err != nil {
		t.Error(err)
	}
	query = "UPDATE `account` SET `is_valid`=? WHERE `id`=?;"
	aff, err := db.Update(query, true, accs[0].ID)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(aff)
}

func TestStruct_Nest(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
		return
	}
	type Contact struct {
		Email  string `sql:"email" json:"email"`
		Mobile string `sql:"mobile" json:"mobile"`
	}
	type Base struct {
		Contact  Contact    `json:"contact"`
		Nickname string     `sql:"nickname" json:"nickname"`
		Avatar   NullString `sql:"avatar" json:"avatar"`
	}
	type Acc struct {
		CreateTime time.Time `sql:"create_time" json:"create_time"`
		ID         int64     `sql:"id" json:"id"`
		Role       NullInt32 `sql:"role" json:"role"`
		Base       Base      `json:"base"`
		Password   string    `sql:"password" json:"password"`
		IsValid    NullBool  `sql:"is_valid" json:"is_valid"`
	}
	var accs []*Acc
	query := "SELECT `id`, `avatar`, `email`, `mobile`, `nickname`, `password`, `role`, `create_time`, `is_valid` FROM `account`;"
	err = db.Query(&accs, query)
	if err != nil {
		t.Error(err)
	}
	resStr, _ := json.Marshal(accs)
	fmt.Println(string(resStr))
	if err := db.Close(); err != nil {
		t.Error(err)
	}
}

func TestStructNest2(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
		return
	}
	type Contact struct {
		Email  string `sql:"email" json:"email"`
		Mobile string `sql:"mobile" json:"mobile"`
	}
	type Base struct {
		Contact  Contact    `json:"contact"`
		Nickname string     `sql:"nickname" json:"nickname"`
		Avatar   NullString `sql:"avatar" json:"avatar"`
	}
	type Acc struct {
		Base     Base   `json:"base"`
		Password string `sql:"password" json:"password"`
	}
	var accs []*Acc
	query := "SELECT  `email`, `avatar`, `mobile`, `nickname`, `password`  FROM `account`;"
	err = db.Query(&accs, query)
	if err != nil {
		t.Error(err)
	}
	resStr, _ := json.Marshal(accs)
	fmt.Println(string(resStr))
	if err := db.Close(); err != nil {
		t.Error(err)
	}
}

func TestStructNest3(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
		return
	}
	type Contact struct {
		Email  string `sql:"email" json:"email"`
		Mobile string `sql:"mobile" json:"mobile"`
	}
	type Base struct {
		Contact  *Contact   `json:"contact"`
		Nickname string     `sql:"nickname" json:"nickname"`
		Avatar   NullString `sql:"avatar" json:"avatar"`
	}
	type Acc struct {
		Base     *Base  `json:"base"`
		Password string `sql:"password" json:"password"`
	}
	var accs []*Acc
	query := "SELECT  `email`, `avatar`, `mobile`, `nickname`, `password`  FROM `account`;"
	err = db.Query(&accs, query)
	if err != nil {
		t.Error(err)
		return
	}
	resStr, _ := json.Marshal(accs)
	fmt.Println(string(resStr))
	if err := db.Close(); err != nil {
		t.Error(err)
	}
}

func TestSqlY_UpdateMany(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
		return
	}

	type accIDs struct {
		ID int64 `json:"id" sql:"id"`
	}
	var ids []*accIDs
	query := "SELECT `id` FROM `account` WHERE `id`<3"
	err = db.Query(&ids, query)
	if err != nil {
		t.Error(err)
		return
	}
	query = "UPDATE `account` SET `password`=? WHERE `id`=?"
	var params [][]interface{}
	for _, id := range ids {
		hash := sha1.New()
		_, _ = hash.Write([]byte(strconv.FormatInt(id.ID, 10)))
		passwd := hex.EncodeToString(hash.Sum(nil))
		params = append(params, []interface{}{passwd, id.ID})
	}

	_, err = db.UpdateMany(query, params)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSqlY_NullTime(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "INSERT INTO `account` (`nickname`, `mobile`, `email`, `add_time`, `birthday`) " +
		"VALUES (?, ?, ?, ?, ?);"
	var vals = [][]interface{}{
		{"💩	'ok'", "18112362345", "testq1@foxmail.com", NullTime{}, NullTime{}},
		{"testq2,\\2", "18112362346", "testq2@foxmail.com", NullTime{Time: time.Now(), Valid: true}, NullTime{Time: time.Now(), Valid: true}},
		{"twt\nwafe", "18112362347", "testq1@foxmail.com", NullTime{}, NullTime{}},
		{"t\\estq4", "18112362348", "testq2@foxmail.com", NullTime{}, NullTime{Time: time.Now(), Valid: true}},
	}
	aff, err := db.InsertMany(query, vals)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(aff)
}

func TestSqlY_BaseType(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}

	query := "SELECT `id` FROM `account` ORDER BY `id`;"
	var vals []int64

	err = db.Query(&vals, query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(vals)
}

func TestSqlY_BaseType2(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}

	query := "SELECT `add_time` FROM `account` ORDER BY `id`;"
	var vals []*NullTime

	err = db.Query(&vals, query)
	if err != nil {
		t.Error(err)
	}
	res, _ := json.Marshal(vals)
	fmt.Println(string(res))
}

func TestSqly_BaseType3(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}

	query := "SELECT `nickname` FROM `account` ORDER BY `id`;"
	var vals []string
	err = db.Query(&vals, query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(vals)
}

func TestSqlY_OneBase(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "SELECT COUNT(*) FROM `account`;"
	var num int
	err = db.Get(&num, query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("num", num)
}

func TestSqlY_OneBase2(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "SELECT `create_time` FROM `account` limit 1;"
	create := &NullTime{}
	err = db.Get(create, query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("create", create)
}

func TestSqly_OneBase3(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "SELECT `nickname` FROM `account` limit 1;"
	var nickname string
	err = db.Get(&nickname, query)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("nickname", nickname)
}

func TestSqlY_QueryOneMap(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}

	acc := make(map[string]interface{})
	query := "SELECT * FROM `account` WHERE `id`=?;"
	err = db.Get(&acc, query, 1)
	if err != nil {
		t.Error(err)
	}
	accStr, _ := json.Marshal(acc)
	fmt.Printf("rows %s", accStr)
}

func TestSqlY_QueryMap(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	var accs []map[string]interface{}
	query := "SELECT * FROM `account`;"

	err = db.Query(&accs, query, nil)
	if err != nil {
		t.Error(err)
	}
	accStr, _ := json.Marshal(accs)
	fmt.Printf("rows %s", accStr)
}

func TestSqly_EmptyArray(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	var accs []map[string]interface{}
	query := "SELECT * FROM `account` WHERE `id` IN ?;"
	var ids []int64
	err = db.Query(&accs, query, ids)
	if err != nil {
		t.Error(err)
	}
	accStr, _ := json.Marshal(accs)
	fmt.Printf("rows %s", accStr)
}

func TestSqly_EmptyArray2(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	query := "UPDATE `account` SET `nickname`=`nickname`+'t' WHERE `id` IN ?;"
	var ids []int64
	aff, err := db.Update(query, ids)
	if !errors.Is(err, ErrEmptyArrayInStatement) {
		t.Error(err)
	}
	fmt.Sprintln(aff)
}

func TestSqlY_Json(t *testing.T) {
	acc := Account{
		ID:         1,
		Nickname:   "nickname",
		Avatar:     NullString{String: "", Valid: false},
		Email:      "123@gmail.com",
		Mobile:     "",
		Role:       NullInt32{Int32: 1, Valid: true},
		Password:   "",
		IsValid:    NullBool{Bool: true, Valid: true},
		CreateTime: time.Now(),
		AddTime:    NullTime{Time: time.Now(), Valid: true},
		Birthday:   NullTime{},
	}
	b, err := json.Marshal(acc)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}

func TestSqlY_Query2(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	// user model
	type Acc struct {
		ID     int64  `sql:"id" json:"id"`
		Email  string `sql:"email" json:"email"`
		Mobile string `sql:"mobile" json:"mobile"`
		Ext    map[string]interface{}
	}
	query := "SELECT * FROM `account`;"
	var accs []*Acc
	err = db.Query(&accs, query)
	if err != nil {
		t.Error(err)
	}
	accStr, _ := json.Marshal(accs)
	fmt.Printf("rows %s", accStr)
}

func TestSqlY_Nest4(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	type Contact struct {
		Email  string                 `sql:"email" json:"email"`
		Mobile string                 `json:"mobile"`
		Ext    map[string]interface{} `sql:"mobile"`
	}
	type Base struct {
		Contact  *Contact   `json:"contact"`
		Nickname string     `sql:"nickname" json:"nickname"`
		Avatar   NullString `sql:"avatar" json:"avatar"`
	}
	type Acc struct {
		Base     *Base  `json:"base"`
		Password string `sql:"password" json:"password"`
	}
	query := "SELECT * FROM `account`;"
	var accs []*Acc
	err = db.Query(&accs, query)
	if err != nil {
		t.Error(err)
	}
	accStr, _ := json.Marshal(accs)
	fmt.Printf("rows %s", accStr)
}
