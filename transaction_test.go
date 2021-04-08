package sqly

import (
	"context"
	"fmt"
	"testing"
)

func TestSqlY_Transaction(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	res, err := db.Transaction(func(tx *Trans) (i interface{}, e error) {
		ctx := context.TODO()
		acc := new(Account)
		query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
			"FROM `account` WHERE `mobile`=?;"
		err := tx.GetCtx(ctx, acc, query, "18812311232")
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `nickname`=? WHERE `id`=?"
		aff, err := tx.Update(query, "nick_trans", acc.ID)
		if err != nil {
			return nil, err
		}
		return aff, nil
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintln(res)
}

func TestSqlY_Transaction2(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	res, err := db.Transaction(func(tx *Trans) (i interface{}, e error) {
		ctx := context.TODO()
		acc := new(Account)
		query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
			"FROM `account` WHERE `mobile`=?;"
		err := tx.GetCtx(ctx, acc, query, "18812311232")
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `nickname`=? WHERE `id`=?"
		aff, err := tx.UpdateCtx(ctx, query, "nick_trans_failed", acc.ID)
		if err != nil {
			return nil, err
		}
		return aff, nil
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintln(res)
}

func TestSqlY_Transaction3(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	res, err := db.Transaction(func(tx *Trans) (i interface{}, e error) {
		ctx := context.TODO()
		query := "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
			"VALUES (?, ?, ?, ?);"
		var vals = [][]interface{}{
			{"testt1", "18112342355", "testq1@foxmail.com", 1},
			{"testt2", "18112342356", "testq2@foxmail.com", 1},
		}
		aff, err := tx.InsertManyCtx(ctx, query, vals)
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `nickname`=? WHERE id=?"
		lastId, err := aff.GetLastId()
		if err != nil {
			return nil, err
		}
		aff, err = tx.UpdateCtx(ctx, query, "last_nick", lastId)
		if err != nil {
			return nil, err
		}
		rows, err := aff.GetRowsAffected()
		if err != nil {
			return nil, err
		}
		lastId, err = aff.GetLastId()
		if err != nil {
			return nil, err
		}
		t.Logf("lastId: %d, rows affected: %d", lastId, rows)
		return aff, nil
	})

	if err != nil {
		t.Error(err)
	}
	fmt.Sprintln(res)
}

func TestSqlY_Transaction4(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	res, err := db.Transaction(func(tx *Trans) (i interface{}, e error) {
		ctx := context.TODO()
		var queries []string
		var mobiles = []string{"18112342355", "18112342356"}
		for _, m := range mobiles {
			item, err := QueryFmt("DELETE FROM `account` WHERE `mobile`=?;", m)
			if err != nil {
				return nil, err
			}
			queries = append(queries, item)
		}
		err := tx.ExecManyCtx(ctx, queries)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintln(res)
}
