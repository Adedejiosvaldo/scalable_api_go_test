package db

import (
	"context"
	"time"
)

type OperationMetrics struct {
	OperationType  string        `json:"operation_type"`
	Duration       time.Duration `json:"duration"`
	BytesProcessed int64         `json:"bytes_processed"`
	CacheHit       bool          `json:"cache_hit"`
	DatabaseType   string        `json:"database_type"`
	IsRemote       bool          `json:"is_remote"`
}

type DatabaseOperation interface {
	Write(ctx context.Context, data interface{}) (*OperationMetrics, error)
	Read(ctx context.Context, id string) (interface{}, *OperationMetrics, error)
	BatchWrite(ctx context.Context, data []interface{}) (*OperationMetrics, error)
	BatchRead(ctx context.Context, ids []string) ([]interface{}, *OperationMetrics, error)
	Benchmark(ctx context.Context, operationType string) (*OperationMetrics, error)
}

type DataStore struct {
	ID        string      `json:"id" bson:"_id"`
	Data      interface{} `json:"data" bson:"data"`
	Timestamp time.Time   `json:"timestamp" bson:"timestamp"`
}
