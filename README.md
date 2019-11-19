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
```go
type Option struct {
	Dsn             string        `json:"dsn"`                // database server name
	DriverName      string        `json:"driver_name"`        // database driver
	MaxIdleConns    int           `json:"max_idle_conns"`     // limit the number of idle connections
	MaxOpenConns    int           `json:"max_open_conns"`     // limit the number of total open connections
	ConnMaxLifeTime time.Duration `json:"conn_max_life_time"` // maximum amount of time a connection may be reused
}
```

> Dsn: 格式化的数据库服务访问参数 例如：[mysql](https://github.com/go-sql-driver/mysql) 格式化方式如下 [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]

> DriverName: 使用的数据库驱动类型 例如： mysql, postgres, sqlite3 等

> MaxIdleConns: 最大空闲连接数

> MaxOpenConns: 最大连接池大小

> ConnMaxLifeTime: 连接的生命周期


详细配置课查看 【Go database/sql tutorial](http://go-database-sql.org/connection-pool.html), [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) 等。

> 连接建立 

