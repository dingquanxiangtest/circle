package dao

import (
	"context"
	"errors"

	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"
)

var (
	// ErrAssertBuilder assert builder fail
	ErrAssertBuilder = errors.New("assert builder fail")
)

// PageData entities of table
type PageData struct {
	Value []Data
	Total int64
}

// FindOptions page options
type FindOptions struct {
	Page int64
	Size int64
	Sort []string
}

// Data entity of table
type Data map[string]interface{}

// Dao dao
type Dao interface {
	FindOne(ctx context.Context, builder clause.Builder) (Data, error)
	Find(ctx context.Context, builder clause.Builder, findOpt FindOptions) ([]Data, error)
	Count(ctx context.Context, builder clause.Builder) (int64, error)
	Insert(ctx context.Context, entity ...interface{}) error
	Update(ctx context.Context, builder clause.Builder, entity interface{}) (int64, error)
	Delete(ctx context.Context, builder clause.Builder) (int64, error)
}
