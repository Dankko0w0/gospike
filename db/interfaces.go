package db

import (
	"context"
)

// Config 数据库配置结构
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Options  map[string]interface{}
}

// DBInterface 定义所有数据库操作的通用接口
type DBInterface interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context) error
	IsConnected() bool
	Reconnect(ctx context.Context) error
}

// DataOperation 定义基础的数据操作接口
type DataOperation interface {
	Create(ctx context.Context, collection string, data interface{}) error
	Read(ctx context.Context, collection string, filter interface{}, result interface{}) error
	Update(ctx context.Context, collection string, filter interface{}, update interface{}) error
	Delete(ctx context.Context, collection string, filter interface{}) error
	List(ctx context.Context, collection string, filter interface{}, results interface{}) error
}
