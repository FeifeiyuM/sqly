package sqly

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestCapsule_InsertUpdate(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	capsule := NewCapsule(db)
	ctx := context.TODO()

	ret, err := capsule.StartCapsule(ctx, true, func(ctx context.Context) (interface{}, error) {
		query := "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
			"VALUES (?, ?, ?, ?);"
		_, err := capsule.Insert(ctx, query, "nick1", "18312311235", "testc1@foxmail.com", 1)
		if err != nil {
			return nil, err
		}
		var params [][]interface{}
		params = append(params, []interface{}{"nick2", "18312311234", "testc2@foxmail.com", 2})
		params = append(params, []interface{}{"nick3", "18312311233", "testc3@foxmail.com", 2})
		params = append(params, []interface{}{"nick4", "18312311232", "testc4@foxmail.com", 1})
		params = append(params, []interface{}{"nick5", "18312311231", "testc5@foxmail.com", 2})
		_, err = capsule.InsertMany(ctx, query, params)
		if err != nil {
			return nil, err
		}

		query = "UPDATE `account` SET `nickname`=? WHERE `mobile`=?"
		var params2 [][]interface{}
		params2 = append(params2, []interface{}{"nick_cap1", "18312311234"})
		params2 = append(params2, []interface{}{"nick_cap2", "18312311233"})
		_, err = capsule.UpdateMany(ctx, query, params2)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintln(ret)
}

func TestCapsule_Delete(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	capsule := NewCapsule(db)
	ctx := context.TODO()

	ret, err := capsule.StartCapsule(ctx, false, func(ctx context.Context) (interface{}, error) {
		query := "SELECT * FROM `account` WHERE `mobile`=?"
		acc := &Account{}
		err := capsule.Get(ctx, acc, query, "18312311231")
		if err != nil {
			return nil, err
		}
		query = "DELETE FROM `account` WHERE `id` IN ?"
		aff, err := capsule.Delete(ctx, query, []int64{acc.ID})
		if err != nil {
			return nil, err
		}
		return aff, err
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintln(ret)
}

func TestCapsule_trans(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	capsule := NewCapsule(db)
	ctx := context.TODO()
	ret, err := capsule.StartCapsule(ctx, true, func(ctx context.Context) (interface{}, error) {
		var accs []*Account
		query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
			"FROM `account`"
		err := capsule.Query(ctx, &accs, query)
		if err != nil {
			return nil, err
		}
		if len(accs) < 2 {
			return nil, nil
		}
		query = "UPDATE `account` SET `nickname`=? WHERE `id`=?"
		_, err = capsule.Update(ctx, query, "nick_trans", accs[0].ID)
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `avatar`=? WHERE `id`=?"
		aff, err := capsule.Update(ctx, query, "test.png", accs[1].ID)
		if err != nil {
			return nil, err
		}
		return aff, nil
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintln(ret)
}

func TestCapsule_trans2(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	capsule := NewCapsule(db)
	ctx := context.TODO()
	ret, err := capsule.StartCapsule(ctx, true, func(ctx context.Context) (interface{}, error) {
		var accs []*Account
		query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
			"FROM `account`"
		err := capsule.Query(ctx, &accs, query, "18812311232")
		if err != nil {
			return nil, err
		}
		if len(accs) < 2 {
			return nil, nil
		}
		query = "UPDATE `account` SET `nickname`=? WHERE `id`=?"
		_, err = capsule.Update(ctx, query, "nick_trans2", accs[0].ID)
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `avatar`=? WHERE `id`=?"
		aff, err := capsule.Update(ctx, query, "test2.png", accs[1].ID)
		if err != nil {
			return nil, err
		}
		query = "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
			"VALUES (?, ?, ?, ?);"
		aff, err = capsule.Insert(ctx, query, "nick_test2", "18712311235", "testx1@foxmail.com", 1)
		if err != nil {
			t.Error(err)
		}
		if aff != nil {
			return nil, errors.New("error")
		}
		return aff, nil
	})
	if err != nil && err.Error() != "error" {
		t.Error(err)
	}
	fmt.Sprintln(ret)
}

func TestCapsule_raw1(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	capsule := NewCapsule(db)
	ctx := context.TODO()
	ret, err := capsule.StartCapsule(ctx, false, func(ctx context.Context) (interface{}, error) {
		var accs []*Account
		query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
			"FROM `account`"
		err := capsule.Query(ctx, &accs, query, "18812311232")
		if err != nil {
			return nil, err
		}
		if len(accs) < 2 {
			return nil, nil
		}
		query = "UPDATE `account` SET `nickname`=? WHERE `id`=?"
		_, err = capsule.Update(ctx, query, "nick_trans3", accs[0].ID)
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `avatar`=? WHERE `id`=?"
		aff, err := capsule.Update(ctx, query, "test3.png", accs[1].ID)
		if err != nil {
			return nil, err
		}
		query = "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
			"VALUES (?, ?, ?, ?);"
		aff, err = capsule.Insert(ctx, query, "nick_test3", "18712311235", "testx1@foxmail.com", 1)
		if err != nil {
			t.Error(err)
		}
		if aff != nil {
			return nil, errors.New("error")
		}
		return aff, nil
	})
	if err != nil && err.Error() != "error" {
		t.Error(err)
	}
	fmt.Sprintln(ret)
}

func TestCapsule_raw2(t *testing.T) {
	db, err := New(opt)
	if err != nil {
		t.Error(err)
	}
	capsule := NewCapsule(db)
	ctx := context.TODO()
	ret, err := capsule.StartCapsule(ctx, false, func(ctx context.Context) (interface{}, error) {
		var accs []*Account
		query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
			"FROM `account`"
		err := capsule.Query(ctx, &accs, query)
		if err != nil {
			return nil, err
		}
		if len(accs) < 2 {
			return nil, nil
		}
		query = "UPDATE `account` SET `nickname`=? WHERE `id`=?"
		_, err = capsule.Update(ctx, query, "nick_trans4", accs[0].ID)
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `avatar`=? WHERE `id`=?"
		aff, err := capsule.Update(ctx, query, "test4.png", accs[1].ID)
		if err != nil {
			return nil, err
		}
		query = "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
			"VALUES (?, ?, ?, ?);"
		aff, err = capsule.Insert(ctx, query, "nick_test4", "18712311230", "testx2@foxmail.com", 1)
		if err != nil {
			t.Error(err)
		}
		return aff, nil
	})
	if err != nil {
		t.Error(err)
	}
	fmt.Sprintln(ret)
}
