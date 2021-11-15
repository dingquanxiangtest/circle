package comet

import "github.com/gin-gonic/gin"

// SchemaHandle SchemaHandle
type SchemaHandle interface {
	SchemaHandle(c *gin.Context, bus *bus) (Pack, error)
}

// DataHandle DataHandle
type DataHandle interface {
	DataHandle(c *gin.Context, bus *bus) (Pack, error)
}

// Search Search
type Search interface {
	Search(c *gin.Context, bus *bus) (Pack, error)
}

// Create CreateData
type Create interface {
	Create(c *gin.Context, bus *bus) (Pack, error)
}

// Update Update
type Update interface {
	Update(c *gin.Context, bus *bus) (Pack, error)
}

// Delete Delete
type Delete interface {
	Delete(c *gin.Context, bus *bus) (Pack, error)
}

// Pre Pre
type Pre interface {
	Pre(c *gin.Context, bus *bus, method string, entity Entity, opts ...PreOption) error
}

// Post Post
type Post interface {
	Postfix(c *gin.Context, bus *bus, pack Pack, opt ...FilterOption) error
}
