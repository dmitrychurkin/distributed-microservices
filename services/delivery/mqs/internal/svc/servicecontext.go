package svc

import (
	"delivery/mqs/internal/config"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	sqlx.SqlConn
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewSqlConn("postgres", c.DataSource)

	return &ServiceContext{
		Config:  c,
		SqlConn: conn,
	}
}
