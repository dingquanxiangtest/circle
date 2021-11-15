package comet

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	code         = "code"
	msg          = "msg"
	data         = "data"
	entities     = "entities"
	total        = "total"
	errorCount   = "errorCount"
	defaultCode  = -1
	entity       = "entity"
	aggregations = "aggregations"
)

// Pack response body pack
type Pack interface {
	Format() map[string]interface{}
	Get() interface{}
}

// Schema JSON schema
type Schema interface{}

// Body response body
type Body struct {
	Data   interface{} `json:"data"`
	Number int64       `json:"number"`
}

// Get Get
func (b *Body) Get() interface{} {
	return b.Data
}

// Format body format
func (b *Body) Format() map[string]interface{} {

	return map[string]interface{}{
		data: map[string]interface{}{
			entity:     b.Data,
			errorCount: b.Number,
		},
	}

}

// Schemas Schemas
type Schemas struct {
	ID      string      `json:"id"`
	TableID string      `json:"tableID"`
	Schema  interface{} `json:"schema"`
	Config  interface{} `json:"config"`
}

// Get get data
func (s *Schemas) Get() interface{} {
	return s.Schema
}

// Format Format
func (s *Schemas) Format() map[string]interface{} {
	return map[string]interface{}{
		data: map[string]interface{}{
			"id":      s.ID,
			"tableID": s.TableID,
			"schema":  s.Schema,
			"config":  s.Config,
		},
	}
}

// Paging paging result set
type Paging struct {
	Data         interface{} `json:"data"`
	Total        int64       `json:"total"`
	Aggregations interface{} `json:"aggregations"`
}

// Format paging format
func (p *Paging) Format() map[string]interface{} {
	return map[string]interface{}{
		data: map[string]interface{}{
			entities:     p.Data,
			total:        p.Total,
			aggregations: p.Aggregations,
		},
	}
}

// Get get data
func (p *Paging) Get() interface{} {
	return p.Data
}

// Resp response
type Resp map[string]interface{}

type packOpt func(Resp)

// WithPack response with pack
func WithPack(pack Pack) func(Resp) {
	return func(resp Resp) {
		if pack == nil {
			return
		}
		for key, value := range pack.Format() {
			resp[key] = value
		}
	}
}

// WithSchema response with schema
func WithSchema(sm Schema) func(Resp) {
	return func(resp Resp) {
		if sm == nil {
			return
		}
		resp[data] = sm
	}
}

// WithError response with error
func WithError(err error) func(Resp) {
	return func(resp Resp) {
		resp[code] = defaultCode
		resp[msg] = err.Error()
	}
}

// Format format response
func Format(c *gin.Context, opts ...packOpt) {
	resp := make(map[string]interface{})
	resp[code] = 0

	for _, opt := range opts {
		opt(resp)
	}
	c.JSON(http.StatusOK, resp)
}
