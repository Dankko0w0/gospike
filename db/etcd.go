package db

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Etcd struct {
	config     *Config
	client     *clientv3.Client
	connected  bool
	maxRetries int
}

func NewEtcd(config *Config) *Etcd {
	return &Etcd{
		config:     config,
		maxRetries: 3,
	}
}

func (e *Etcd) Connect(ctx context.Context) error {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)},
		Username:    e.config.Username,
		Password:    e.config.Password,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return fmt.Errorf("failed to connect to etcd: %v", err)
	}

	e.client = client
	e.connected = true
	return nil
}

func (e *Etcd) Disconnect(ctx context.Context) error {
	if e.client != nil {
		err := e.client.Close()
		if err != nil {
			return fmt.Errorf("failed to disconnect from etcd: %v", err)
		}
		e.connected = false
	}
	return nil
}

func (e *Etcd) Ping(ctx context.Context) error {
	if e.client == nil {
		return fmt.Errorf("database not connected")
	}

	_, err := e.client.Get(ctx, "ping")
	return err
}

func (e *Etcd) IsConnected() bool {
	return e.connected
}

func (e *Etcd) Reconnect(ctx context.Context) error {
	e.Disconnect(ctx)

	for i := 0; i < e.maxRetries; i++ {
		err := e.Connect(ctx)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return fmt.Errorf("failed to reconnect after %d attempts", e.maxRetries)
}

// Etcd specific operations
func (e *Etcd) Put(ctx context.Context, key, value string) error {
	_, err := e.client.Put(ctx, key, value)
	return err
}

func (e *Etcd) Get(ctx context.Context, key string) (string, error) {
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("key not found")
	}
	return string(resp.Kvs[0].Value), nil
}

func (e *Etcd) Delete(ctx context.Context, key string) error {
	_, err := e.client.Delete(ctx, key)
	return err
}

func (e *Etcd) Watch(ctx context.Context, key string) clientv3.WatchChan {
	return e.client.Watch(ctx, key)
}
