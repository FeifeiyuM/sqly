package sqlyt

import (
	"context"
	"errors"
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
			"FROM `accounts` WHERE `mobile`=?;"
		err := tx.QueryOneCtx(ctx, acc, query, "18812311232")
		if err != nil {
			return nil, err
		}
		query = "UPDATE `accounts` SET `nickname`=? WHERE `id`=?"
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
			"FROM `accounts` WHERE `mobile`=?;"
		err := tx.QueryOneCtx(ctx, acc, query, "18787655678")
		if err != nil {
			return nil, err
		}
		query = "UPDATE `accounts` SET `nickname`= WHERE `id`=?"
		aff, err := tx.UpdateCtx(ctx, query, "nick_trans_failed", acc.ID)
		if err == nil {
			return nil, errors.New("new error")
		}
		return aff, nil
	})
	if err == nil {
		t.Error("need err not null")
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
		query := "INSERT INTO `accounts` (`nickname`, `mobile`, `email`, `role`) " +
			"VALUES (?, ?, ?, ?);"
		var vals = [][]interface{}{
			{"testt1", "18112342355", "testq1@foxmail.com", 1},
			{"testt2", "18112342356", "testq2@foxmail.com", 1},
		}
		aff, err := tx.InsertManyCtx(ctx, query, vals)
		if err != nil {
			return nil, err
		}
		query = "UPDATE `accounts` SET `nickname`=? WHERE id=?"
		aff, err = tx.UpdateCtx(ctx, query, "last_nick", aff.LastId)
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

func TestSqlY_Transaction4(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	res, err := db.Transaction(func(tx *Trans) (i interface{}, e error) {
		ctx := context.TODO()
		var quries []string
		var mobiles = []string{"18112342355", "18112342356"}
		for _, m := range mobiles {
			item, err := QueryFmt("DELETE FROM `accounts` WHERE `mobile`=?;", m)
			if err != nil {
				return nil, err
			}
			quries = append(quries, item)
		}
		err := tx.ExecManyCtx(ctx, quries)
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
