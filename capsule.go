package sqly

import "context"

// Capsule 胶囊对象
type Capsule struct {
	sqlY *SqlY
}

// CapFunc 胶囊闭包函数
type CapFunc func(ctx context.Context) (interface{}, error)

// NewCapsule new capsule
func NewCapsule(sqlY *SqlY) *Capsule {
	return &Capsule{sqlY: sqlY}
}

// Capsule 事务胶囊
type capsule struct {
	tx      *Trans // 事务句柄
	conn    *SqlY  // 普通连接
	isTrans bool   // 是否开启事务
}

// GetCapsule 获取连接胶囊
func (c *Capsule) getCapsule(ctx context.Context) (*capsule, error) {
	ci := ctx.Value("_sqly_capsule")
	if ci == nil {
		return &capsule{conn: c.sqlY}, nil
	}
	cs, ok := ci.(*capsule)
	if !ok {
		return &capsule{
			tx:      nil,
			conn:    c.sqlY,
			isTrans: false,
		}, nil
	}
	if cs.isTrans {
		if cs.tx == nil {
			return nil, ErrCapsule
		}
	} else {
		if cs.conn == nil {
			return nil, ErrCapsule
		}
	}
	return cs, nil
}

// StartCapsule 开启查询胶囊
func (c *Capsule) StartCapsule(ctx context.Context, isTrans bool, capFunc CapFunc) (interface{}, error) {
	capsule := &capsule{isTrans: isTrans}
	if !isTrans {
		capsule.conn = c.sqlY
		sCtx := context.WithValue(ctx, "_sqly_capsule", capsule)
		return capFunc(sCtx)
	} else {
		return c.sqlY.Transaction(func(tx *Trans) (interface{}, error) {
			capsule.tx = tx
			sCtx := context.WithValue(ctx, "_sqly_capsule", capsule)
			return capFunc(sCtx)
		})
	}
}

// Query query
func (c *Capsule) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	cs, err := c.getCapsule(ctx)
	if err != nil {
		return err
	}
	if cs.isTrans {
		return cs.tx.QueryCtx(ctx, dest, query, args...)
	}
	return cs.conn.QueryCtx(ctx, dest, query, args...)
}

// Get query one
func (c *Capsule) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	cs, err := c.getCapsule(ctx)
	if err != nil {
		return err
	}
	if cs.isTrans {
		return cs.tx.GetCtx(ctx, dest, query, args...)
	}
	return cs.conn.GetCtx(ctx, dest, query, args...)
}

// Insert insert
func (c *Capsule) Insert(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	cs, err := c.getCapsule(ctx)
	if err != nil {
		return nil, err
	}
	if cs.isTrans {
		return cs.tx.InsertCtx(ctx, query, args...)
	} else {
		return cs.conn.InsertCtx(ctx, query, args...)
	}
}

// InsertMany insert many
func (c *Capsule) InsertMany(ctx context.Context, query string, args [][]interface{}) (*Affected, error) {
	cs, err := c.getCapsule(ctx)
	if err != nil {
		return nil, err
	}
	if cs.isTrans {
		return cs.tx.InsertManyCtx(ctx, query, args)
	} else {
		return cs.conn.InsertCtx(ctx, query, args)
	}
}

// Update update
func (c *Capsule) Update(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	cs, err := c.getCapsule(ctx)
	if err != nil {
		return nil, err
	}
	if cs.isTrans {
		return cs.tx.UpdateCtx(ctx, query, args...)
	} else {
		return cs.conn.UpdateCtx(ctx, query, args...)
	}
}

// UpdateMany update many
func (c *Capsule) UpdateMany(ctx context.Context, query string, args [][]interface{}) (*Affected, error) {
	cs, err := c.getCapsule(ctx)
	if err != nil {
		return nil, err
	}
	if cs.isTrans {
		return cs.tx.UpdateManyCtx(ctx, query, args)
	} else {
		return cs.conn.UpdateManyCtx(ctx, query, args)
	}
}

// Delete delete
func (c *Capsule) Delete(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	cs, err := c.getCapsule(ctx)
	if err != nil {
		return nil, err
	}
	if cs.isTrans {
		return cs.tx.DeleteCtx(ctx, query, args...)
	} else {
		return cs.conn.DeleteCtx(ctx, query, args...)
	}
}

// Exec exec
func (c *Capsule) Exec(ctx context.Context, query string, args ...interface{}) (*Affected, error) {
	cs, err := c.getCapsule(ctx)
	if err != nil {
		return nil, err
	}
	if cs.isTrans {
		return cs.tx.ExecCtx(ctx, query, args...)
	} else {
		return cs.conn.ExecCtx(ctx, query, args...)
	}
}

// ExecMany exec multi queries
func (c *Capsule) ExecMany(ctx context.Context, queries []string) error {
	cs, err := c.getCapsule(ctx)
	if err != nil {
		return err
	}
	if cs.isTrans {
		return cs.tx.ExecManyCtx(ctx, queries)
	} else {
		return cs.conn.ExecManyCtx(ctx, queries)
	}
}
