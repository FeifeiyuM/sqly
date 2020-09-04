# sqly

sqly 是基于 golang s数据库操作的标准包 database/sql 的扩展。

[![Build Status](https://action-badges.now.sh/FeifeiyuM/sqly)](https://github.com/FeifeiyuM/sqly/actions?query=workflow%3AGo)
[![Go Report](https://goreportcard.com/badge/github.com/FeifeiyuM/sqly)](https://goreportcard.com/report/github.com/FeifeiyuM/sqly)
[![Coverage Status](https://coveralls.io/repos/github/FeifeiyuM/sqly/badge.svg?branch=master)](https://coveralls.io/github/FeifeiyuM/sqly?branch=master)

主要目标（功能)：
- 是实现类似于 json.Marshal 类似的功能，将数据库查询结果反射成为 struct 对象。
简化 database/sql 原生的 span 书写方法。

- 通过回调函数的形式封装了事务操作，简化原生包关于事务逻辑的操作

- 封装了原生 database/sql 包不具有的, 更新（Update), 插入(Insert), 删除（DELETE), 通用sql 执行(Exec) 等方法（Exec）


## 使用

### 安装依赖
> go get github.com/FeifeiyuM/sqly

### 连接数据库
> **连接配置**
 func New(opt *sqly.Option) (*SqlY, error)
```go
    opt := &sqly.Option{
		Dsn:             "test:mysql123@tcp(127.0.0.1:3306)/test_db?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local",
		DriverName:      "mysql",
		MaxIdleConns:    0,
		MaxOpenConns:    0,
		ConnMaxLifeTime: 0,
	}
    db, err := sqly.New(opt)
	if err != nil {
		fmt.Println("test error")
	}
    // 数据库连接成功
```

> Dsn: 格式化的数据库服务访问参数 例如：[mysql](https://github.com/go-sql-driver/mysql) 格式化方式如下 [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]

> DriverName: 使用的数据库驱动类型 例如： mysql, postgres, sqlite3 等

> MaxIdleConns: 最大空闲连接数

> MaxOpenConns: 最大连接池大小

> ConnMaxLifeTime: 连接的生命周期


详细配置课查看 【Go database/sql tutorial](http://go-database-sql.org/connection-pool.html), [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) 等。


### 数据库操作
- 通用执行操作, 执行一次命令（包括查询、删除、更新、插入, 建表等）
> func (s *SqlY) Exec(query string, args ...interface{}) (*Affected, error)

> func (s *SqlY) ExecCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)
```go
    // 创建表
    query := "CREATE TABLE `account` (" +
    		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
    		"`nickname` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL," +
    		"`avatar` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'avatar url'," +
    		"`mobile` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'mobile number'," +
    		"`email` varchar(320) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'email'," +
    		"`password` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT 'password'," +
    		"`role` tinyint(4) DEFAULT NULL COMMENT 'role'," +
    		"`expire_time` datetime DEFAULT NULL COMMENT 'expire_time'," +
    		"`is_valid` tinyint(4) DEFAULT NULL COMMENT 'is_valid'," +
    		"`create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP," +
    		"PRIMARY KEY (`id`)," +
    		"UNIQUE KEY `mobile_index` (`mobile`)," +
    		"KEY `email_index` (`email`)" +
    		") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;"
    _, err = db.Exec(query)
    if err != nil {
    	fmt.Println("create table error")
    }

```

- 插入一条数据 
> func (s *SqlY) Insert(query string, args ...interface{}) (*Affected, error) 
> func (s *SqlY) InsertCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)
```go
    query := "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) VALUES (?, ?, ?, ?);"
	aff, err := db.Insert(query, "nick_test3", "18812311235", "test@foxmail.com", 1)
	if err != nil {
		fmt.Println("failed to insert data")
	}
    if aff != nil {
        fmt.Printf("auto_id: %v, affected_rows: %v\n", aff.LastId, aff.RowsAffected)
    }
   // Affected 结构体返回的值，不保证值(LastId, RowsAffected)不为空，根据不同数据库和其对应的驱动确定
   // lastId 表示最后一条插入数据对应有数据生成的一个数字id(自增id), 
   // RowsAffected 表示 update, insert, or delete 操作影响的行数。
```

- 插入多条数据 
> func (s *SqlY) InsertMany(query string, args [][]interface{}) (*Affected, error)

> func (s *SqlY) InsertManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error)
```go
    query := "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) VALUES (?, ?, ?, ?);"
    var vals = [][]interface{}{
        {"testq1", "18112342345", "testq1@foxmail.com", 1},
        {"testq2", "18112342346", "testq2@foxmail.com", nil},
    }
    aff, err = db.InsertMany(query, vals)
    if err != nil {
        fmt.Sprintln("create account error")
    }
    if err != nil {
        fmt.Sprintln("create accounts error")
    }
    fmt.Println(aff)
```

- 更新一条数据 
> func (s *SqlY) Update(query string, args ...interface{}) (*Affected, error)

> func (s *SqlY) ExecCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)
```go
    query := "UPDATE `account` SET `nickname`=? WHERE `mobile`=?;"
	aff, err := db.Update(query, "lucy", "18812311231")
	if err != nil {
		fmt.Sprintln("update accounts error")
	}
	fmt.Println(aff)
```

- 更新多条数据
> func (s *SqlY) UpdateMany(query string, args [][]interface{}) (*Affected, error)
> func (s *SqlY) UpdateManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error)
```go
    query = "UPDATE `account` SET `password`=? WHERE `id`=?"
	var params [][]interface{}
	for _, id := range ids {
		hash := sha1.New()
		_, _ = hash.Write([]byte(strconv.FormatInt(id.ID, 10)))
		passwd := hex.EncodeToString(hash.Sum(nil))
		params = append(params, []interface{}{passwd, id.ID})
	}
	_, err := db.UpdateMany(query, params)
	if err != nil {
		t.Error(err)
		return
	}
```

- 删除一条数据 
> func (s *SqlY) Delete(query string, args ...interface{}) (*Affected, error)

> func (s *SqlY) DeleteCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)
```go
    query := "DELETE FROM `account` WHERE `mobile`=?;"
	aff, err := db.Delete(query, "18812311231")
	if err != nil {
		fmt.Sprintln("delete account error")
	}
	fmt.Println(aff)
```

- 查询一条数据 
> func (s *SqlY) Get(dest interface{}, query string, args ...interface{}) error

> func (s *SqlY) GetCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error
```go
	type Account struct {
		ID         int64      `sql:"id" json:"id"`
		Nickname   string     `sql:"nickname" json:"nickname"`
		Avatar     sqly.NullString `sql:"avatar" json:"avatar"`
		Email      string     `sql:"email" json:"email"`
		Mobile     string     `sql:"mobile" json:"mobile"`
		Role       sqly.NullInt32     `sql:"role" json:"role"`
		Password   string     `sql:"password" json:"password"`
		ExpireTime sqly.NullTime `sql:"expire_time" json:"expire_time"`
		IsValid sqly.NullBool `sql:"is_valid" json:"is_valid"`
		CreateTime time.Time  `sql:"create_time" json:"create_time"`
	}
	acc := new(Account)
	query = "SELECT * FROM `account` WHERE `mobile`=?"
	err = db.Get(acc, query, "18812311235")
	if err != nil {
		fmt.Println("query account error")
	}
	accStr, err := json.Marshal(acc1)
	if err != nil {
		fmt.Println("marshal acc error")
	}
    fmt.Println(accStr)
```
参数 dest 必须为实例化的 struct 对象指针

- 查询数据
> func (s *SqlY) Query(dest interface{}, query string, args ...interface{}) error 

> func (s *SqlY) QueryCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error
```go
    type Account struct {
        ID         int64      `sql:"id" json:"id"`
        Nickname   string     `sql:"nickname" json:"nickname"`
        Avatar     sqly.NullString `sql:"avatar" json:"avatar"`
        Email      string     `sql:"email" json:"email"`
        Mobile     string     `sql:"mobile" json:"mobile"`
        Role       sqly.NullInt32     `sql:"role" json:"role"`
        Password   string     `sql:"password" json:"password"`
        ExpireTime sqly.NullTime `sql:"expire_time" json:"expire_time"`
        IsValid sqly.NullBool `sql:"is_valid" json:"is_valid"`
        CreateTime time.Time  `sql:"create_time" json:"create_time"`
    }

    query = "SELECT * FROM `account` WHERE `mobile` IN ?"
	var mobiles = []string{"18812311235", "18112342346"}
	var accs []*Account  // 必须是 struct array
	err = db.Query(&accs, query, mobiles)
	if err != nil {
		fmt.Printf("query accounts error")
	}
	accsStr, err := json.Marshal(accs)
	if err != nil {
		fmt.Println("marshal acc error")
	}
    fmt.Println(accsStr)
```
参数 dest 必须为实例化的 struct 对象(或对象指针)数组的指针 
     
    
### 数据库事务
- 事务开启
提交，回滚  
> func (s *SqlY) NewTrans() (*Trans, error) 开启

- 事务提交
> func (t *Trans) Commit() error 提交

- 事务回滚
> func (t *Trans) Rollback() error 回滚

- 事务通用执行
> func (t *Trans) Exec(query string, args ...interface{}) (*Affected, error)

> func (t *Trans) ExecCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)

- 事务插入
> func (t *Trans) Insert(query string, args ...interface{}) (*Affected, error)

> func (t *Trans) InsertCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)

- 事务插入多条
> func (t *Trans) InsertMany(query string, args [][]interface{}) (*Affected, error)

> func (t *Trans) InsertManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error)

- 事务更新
> func (t *Trans) Update(query string, args ...interface{}) (*Affected, error)

> func (t *Trans) UpdateCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)

- 事务更新多条
> func (t *Trans) UpdateMany(query string, args [][]interface{}) (*Affected, error)
> func (t *Trans) UpdateManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error)

- 事务删除
> func (t *Trans) Delete(query string, args ...interface{}) (*Affected, error)

> func (t *Trans) DeleteCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)

- 事务查询单条
> func (t *Trans) Get(dest interface{}, query string, args ...interface{}) error 

> func (t *Trans) GetCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error
参数 dest 必须为实例化的 struct 对象指针

- 事务查询
> func (t *Trans) Query(dest interface{}, query string, args ...interface{}) error

> func (t *Trans) QueryCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error
参数 dest 必须为实例化的 struct 对象(或对象指针)数组的指针

```go
    ctx := context.TODO()
    db, err := sqly.New(opt)
	if err != nil {
		fmt.Println("test error")
	}
    // 开始事务
	tx, err := sy.NewTrans()
    if err != nil {
		fmt.Printf("failed to begin transaction")
		return
	}
    // 回滚
	defer func() {
		_ = tx.Rollback()
    }()
    
    type Account struct {
        ID         int64      `sql:"id" json:"id"`
        Nickname   string     `sql:"nickname" json:"nickname"`
        Avatar     sqly.NullString `sql:"avatar" json:"avatar"`
        Email      string     `sql:"email" json:"email"`
        Mobile     string     `sql:"mobile" json:"mobile"`
        Role       sqly.NullInt32     `sql:"role" json:"role"`
        Password   string     `sql:"password" json:"password"`
        ExpireTime sqly.NullTime `sql:"expire_time" json:"expire_time"`
        IsValid sqly.NullBool `sql:"is_valid" json:"is_valid"`
        CreateTime time.Time  `sql:"create_time" json:"create_time"`
    }
    // 执行事务
    // 查
    acc := new(Account)
	query = "SELECT * FROM `account` WHERE `mobile`=?"
	err = tx.GetCtx(ctx, acc, query, "18812311235")
	if err != nil {
		fmt.Printf("get accout error")
		return
	}
    // 更新
    query = "UPDATE `account` SET `is_valid`=? WHERE id=?"
	aff, err := tx.UpdateCtx(ctx, query, true, acc.ID)
	if err != nil {
		fmt.Println("update account error")
	}
    fmt.Println(aff)
    // 删除
    query = "DELETE FROM `account` WHERE id!=?"
    _, err = tx.DeleteCtx(ctx, query, acc.ID)
    if err != nil {
    	fmt.Println("delete accounts error")
    }
    // 插入
    query = "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) VALUES (?, ?, ?, ?);"
	aff, err = tx.InsertCtx(ctx, query, "nick_ruby", "13565656789", nil)
	if err != nil {
		fmt.Println("insert account error")
	}
    fmt.Println(aff)
    // 提交
	_ = tx.Commit()
```

- 事务回调(封装事务开启，关闭，回滚操作)
> type TxFunc func(tx *Trans) (interface{}, error)

> func (s *SqlY) Transaction(txFunc TxFunc) (interface{}, error)
```go
    ctx := context.TODO()
    db, err := sqly.New(opt)
    if err != nil {
    	fmt.Println("test error")
    }
    res, err := db.Transaction(func(tx *sqly.Trans) (i interface{}, e error) {
        // 不需要手动开启，关闭，回滚事务
        // 查
        acc := new(Account)
        query = "SELECT * FROM `account` WHERE `mobile`=?"
        err = tx.GetCtx(ctx, acc, query, "18812311235")
        if err != nil {
            fmt.Printf("get accout error")
            return
        }
        // 更新
        query = "UPDATE `account` SET `is_valid`=? WHERE id=?"
        aff, err := tx.UpdateCtx(ctx, query, true, acc.ID)
        if err != nil {
            fmt.Println("update account error")
        }
        fmt.Println(aff)
        // 删除
        query = "DELETE FROM `account` WHERE id!=?"
        _, err = tx.DeleteCtx(ctx, query, acc.ID)
        if err != nil {
            fmt.Println("delete accounts error")
        }
        // 插入
        query = "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) VALUES (?, ?, ?, ?);"
        aff, err = tx.InsertCtx(ctx, query, "nick_ruby", "13565656789", nil)
        if err != nil {
            fmt.Println("insert account error")
        }
        fmt.Println(aff)
    })
    if err := nil {
        fmt.Println("do transaction error")
    }   
    fmt.Println(res)
```
    
    
### 胶囊查询
> 在执行事务操作的时候，我们需要显式得去初始化和传递事务句柄 tx,  常常不注意就会出现事务和非事务操作混用的问题，严重时还会出现查询线程池耗尽产生死锁（在事务内，申请非事务查询线程，在高并发的时候会尝试该死锁）   
> 为了减少在开发过程中减少对事务和非事务的关注，sqly 采用闭包的方式封装一系列数据库操作，并采用 context 的方式在函数之间传递事务句柄，只需要在初始化闭包的时候确认是否开始事务。

### 胶囊操作相关方法
- 初始化胶囊句柄
> func NewCapsule(sqlY *SqlY) *Capsule 

- 开启胶囊操作
> type CapFunc func(ctx context.Context) (interface{}, error)     
> func (c *Capsule) StartCapsule(ctx context.Context, isTrans bool, capFunc CapFunc) (interface{}, error)     
> StartCapsule 开启胶囊，参数 ctx 上下文用于携带胶囊句柄，isTrans 是否开始事务 true 开启。 CapFunc 闭包函数，所有逻辑都在该闭包内实现

> func (c *Capsule) Exec(query string, args ...interface{}) (*Affected, error)

> func (c *Capsule) ExecCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)

- 插入
> func (c *Capsule) Insert(query string, args ...interface{}) (*Affected, error)

> func (c *Capsule) InsertCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)

- 插入多条
> func (c *Capsule) InsertMany(query string, args [][]interface{}) (*Affected, error)

> func (c *Capsule) InsertManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error)

- 更新
> func (c *Capsule) Update(query string, args ...interface{}) (*Affected, error)

> func (c *Capsule) UpdateCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)

- 更新多条
> func (c *Capsule) UpdateMany(query string, args [][]interface{}) (*Affected, error)
> func (c *Capsule) UpdateManyCtx(ctx context.Context, query string, args [][]interface{}) (*Affected, error)

- 删除
> func (c *Capsule) Delete(query string, args ...interface{}) (*Affected, error)

> func (c *Capsule) DeleteCtx(ctx context.Context, query string, args ...interface{}) (*Affected, error)

- 查询单条
> func (c *Capsule) Get(dest interface{}, query string, args ...interface{}) error 

> func (c *Capsule) GetCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error
参数 dest 必须为实例化的 struct 对象指针

- 查询
> func (c *Capsule) Query(dest interface{}, query string, args ...interface{}) error

> func (c *Capsule) QueryCtx(ctx context.Context, dest interface{}, query string, args ...interface{}) error
参数 dest 必须为实例化的 struct 对象(或对象指针)数组的指针

#### tips: 以上 sql 操作会根据开启胶囊时是否开启事务（isTrans) 设置来自动选择采用事务查询，还是非事务查询。

胶囊例1、事务操作
```go
func TestCapsule_trans2(t *testing.T) {
	db, err := New(opt)  // 初始化sqly(数据库连接）
	if err != nil {
		t.Error(err)
	}
	capsule := NewCapsule(db)  // 创建一个胶囊句柄
	ctx := context.TODO()  // 胶囊查询必须携带 context 参数
    // isTrans=true 开始事务查询
	ret, err := capsule.StartCapsule(ctx, true, func(ctx context.Context) (interface{}, error) {
        // 在闭包内执行相关数据库操作
		var accs []*Account
		query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
			"FROM `account`"
		err := capsule.Query(ctx, &accs, query, "18812311232")  // 执行查询
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `nickname`=? WHERE `id`=?" 
		_, err = capsule.Update(ctx, query, "nick_trans2", accs[0].ID)  // 更新
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `avatar`=? WHERE `id`=?"
		aff, err := capsule.Update(ctx, query, "test2.png", accs[1].ID) // 更新
		if err != nil {
			return nil, err
		}
		query = "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
			"VALUES (?, ?, ?, ?);"
		aff, err = capsule.Insert(ctx, query, "nick_test2", "18712311235", "testx1@foxmail.com", 1)  // 插入
		if err != nil {
			t.Error(err)
		}
		if aff != nil {
			return nil, errors.New("error")
		}
		return aff, nil
	})
	if err.Error() != "error" {
		t.Error(err)
	}
	fmt.Sprintln(ret)
}
```
胶囊例2、非事务操作
```go
func TestCapsule_raw(t *testing.T) {
	db, err := New(opt)  // 初始化sqly(数据库连接）
	if err != nil {
		t.Error(err)
	}
	capsule := NewCapsule(db)  // 创建一个胶囊句柄
	ctx := context.TODO()  // 胶囊查询必须携带 context 参数
    // isTrans=false 不开启事务
	ret, err := capsule.StartCapsule(ctx, false, func(ctx context.Context) (interface{}, error) {
        // 在闭包内执行相关数据库操作
		var accs []*Account
		query := "SELECT `id`, `nickname`, `avatar`, `email`, `mobile`, `password`, `role` " +
			"FROM `account`"
		err := capsule.Query(ctx, &accs, query, "18812311232")  // 查询
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `nickname`=? WHERE `id`=?"
		_, err = capsule.Update(ctx, query, "nick_trans3", accs[0].ID) // 更新
		if err != nil {
			return nil, err
		}
		query = "UPDATE `account` SET `avatar`=? WHERE `id`=?"
		aff, err := capsule.Update(ctx, query, "test3.png", accs[1].ID)  // 更新
		if err != nil {
			return nil, err
		}
		query = "INSERT INTO `account` (`nickname`, `mobile`, `email`, `role`) " +
			"VALUES (?, ?, ?, ?);"
		aff, err = capsule.Insert(ctx, query, "nick_test3", "18712311235", "testx1@foxmail.com", 1)  // 插入
		if err != nil {
			t.Error(err)
		}
		if aff != nil {
			return nil, errors.New("error")
		}
		return aff, nil
	})
	if err.Error() != "error" {
		t.Error(err)
	}
	fmt.Sprintln(ret)
}
```
从胶囊例1和胶囊例2中可以发现，采用胶囊操作时，事务操作和非事务操作区别仅在 StartCapsule 中决定是否开始事务，在实现业务逻辑的查询，更新查询等操作的过程中都无需关注是否开启事务。
      
      
### 支持类型
- struct 中定义的字段类型须是 database/sql 中能够被 Scan 的类型 (int64, float64, bool, []byte, string, time.Time, nil)

- 为了支持更好为空(NULL)的字段，sqly 扩展了 sql.NullTime, sql.NullBool, sql.NullFloat64, sql.NullInt64, sql.NullInt32, 
sql.NullString, 分别为 sqly.NullTime, sqly.NullBool, sqly.NullFloat64, sqly.NullInt64, sqly.NullInt32, sqly.NullString。

- 使用 sqly 扩展的空字段类型，对象在使用 json.Marshal 时 对应字段为空的会自动解析为 null; json 字符串使用 json.UnMarshal 时，会自动解析为对应的 sqly.NullTime 等扩展类型

- 如果使用 tinyint 或 int 类表示 bool 字段类型，例如：0 为 false, 1或**其它**为 true, 在定义字段类型时，可以使用 sqly.Boolean 类型来支持，在 scan 的时候会字段将 int 类型转换成 bool, 如果值只有 0 或 1 可以使用原生 bool

- struct 嵌套支持
```go
    db, err := New(opt)
	if err != nil {
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
		ID         int64     `sql:"id" json:"id"`
		Role       NullInt32 `sql:"role" json:"role"`
		Base       Base      `json:"base"`
		Password   string    `sql:"password" json:"password"`
		IsValid    NullBool  `sql:"is_valid" json:"is_valid"`
        CreateTime time.Time `sql:"create_time" json:"create_time"`
	}
	var accs []*Acc
	query := "SELECT `id`, `avatar`, `email`, `mobile`, `nickname`, `password`, `role`, `create_time`, `is_valid` FROM `account`;"
	err = db.Query(&accs, query)
	if err != nil {
		fmt.Println("query account error")
        reutrn 
	}
	resStr, _ := json.Marshal(accs)
	fmt.Println(string(resStr))
```

- map[string]inteface{} 类型支持(目前支持只 MySQL)
```go
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
```

- 可scan 类型支持 int, int32, int64, string, time.Time, 空字段类型 (及其他们的数组结构)
```go
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
```

```go
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
```

```go
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
```

### tips
- 如果要使用 time.Time 的字段类型, 连接数据库的 dsn 配置中加上 parseTime=true  


