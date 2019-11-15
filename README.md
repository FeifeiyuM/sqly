# sqly

sqly 是基于 golang s数据库操作的标准包 database/sql 的扩展。

主要目标（功能)：
- 是实现类似于 json.Marshal 类似的功能，将数据库查询结果反射成为 struct 对象。
简化 database/sql 原生的 span 书写方法。

- 通过回调函数的形式封装了事务操作，简化原生包关于事务逻辑的操作


## 使用
