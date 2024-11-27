package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	mssql "github.com/microsoft/go-mssqldb"
)

type SQLServer struct {
	config     *Config
	db         *sql.DB
	connected  bool
	maxRetries int
}

func NewSQLServer(config *Config) *SQLServer {
	return &SQLServer{
		config:     config,
		maxRetries: 3,
	}
}

func (s *SQLServer) Connect(ctx context.Context) error {
	// 构建连接字符串
	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		s.config.Username,
		s.config.Password,
		s.config.Host,
		s.config.Port,
		s.config.Database,
	)

	// 创建连接器
	connector, err := mssql.NewConnector(connString)
	if err != nil {
		return fmt.Errorf("failed to create SQL Server connector: %v", err)
	}

	// 创建数据库连接
	db := sql.OpenDB(connector)

	// 设置连接池参数
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	// 测试连接
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return fmt.Errorf("failed to ping SQL Server: %v", err)
	}

	s.db = db
	s.connected = true
	return nil
}

func (s *SQLServer) Disconnect(ctx context.Context) error {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			return fmt.Errorf("failed to disconnect from SQL Server: %v", err)
		}
		s.connected = false
	}
	return nil
}

func (s *SQLServer) Ping(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("database not connected")
	}
	return s.db.PingContext(ctx)
}

func (s *SQLServer) IsConnected() bool {
	return s.connected
}

func (s *SQLServer) Reconnect(ctx context.Context) error {
	s.Disconnect(ctx)

	for i := 0; i < s.maxRetries; i++ {
		err := s.Connect(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return fmt.Errorf("failed to reconnect after %d attempts", s.maxRetries)
}

// CRUD operations
func (s *SQLServer) Create(ctx context.Context, table string, data map[string]interface{}) error {
	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	placeholders := make([]string, 0, len(data))

	i := 0
	for col, val := range data {
		columns = append(columns, col)
		values = append(values, val)
		placeholders = append(placeholders, fmt.Sprintf("@p%d", i))
		i++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		joinStrings(columns, ","),
		joinStrings(placeholders, ","),
	)

	_, err := s.db.ExecContext(ctx, query, values...)
	return err
}

func (s *SQLServer) Read(ctx context.Context, table string, filter map[string]interface{}, result interface{}) error {
	whereClause, values := buildWhereClause(filter)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, whereClause)

	row := s.db.QueryRowContext(ctx, query, values...)
	return row.Scan(result)
}

func (s *SQLServer) Update(ctx context.Context, table string, filter map[string]interface{}, update map[string]interface{}) error {
	setClauses := make([]string, 0, len(update))
	values := make([]interface{}, 0, len(update)+len(filter))

	i := 0
	for col, val := range update {
		setClauses = append(setClauses, fmt.Sprintf("%s = @p%d", col, i))
		values = append(values, val)
		i++
	}

	whereClause, whereValues := buildWhereClause(filter)
	values = append(values, whereValues...)

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table,
		joinStrings(setClauses, ","),
		whereClause,
	)

	_, err := s.db.ExecContext(ctx, query, values...)
	return err
}

func (s *SQLServer) Delete(ctx context.Context, table string, filter map[string]interface{}) error {
	whereClause, values := buildWhereClause(filter)
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", table, whereClause)

	_, err := s.db.ExecContext(ctx, query, values...)
	return err
}

func (s *SQLServer) List(ctx context.Context, table string, filter map[string]interface{}, results interface{}) error {
	whereClause, values := buildWhereClause(filter)
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, whereClause)

	rows, err := s.db.QueryContext(ctx, query, values...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return scanRows(rows, results)
}

// 辅助函数
func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func buildWhereClause(filter map[string]interface{}) (string, []interface{}) {
	if len(filter) == 0 {
		return "1=1", nil
	}

	clauses := make([]string, 0, len(filter))
	values := make([]interface{}, 0, len(filter))

	i := 0
	for col, val := range filter {
		clauses = append(clauses, fmt.Sprintf("%s = @p%d", col, i))
		values = append(values, val)
		i++
	}

	return joinStrings(clauses, " AND "), values
}

// scanRows 扫描结果集到目标结构
func scanRows(rows *sql.Rows, dest interface{}) error {
	// 这里需要根据dest的类型来实现具体的扫描逻辑
	// 可以使用反射来处理不同类型的dest
	return nil
}
