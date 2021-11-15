package restful

import (
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SubTable SubTable
type SubTable struct {
	subTable      service.SubTable
}

// NewSubTable new a manager subTable
func NewSubTable(conf *config.Config, opt ...service.Options) (*SubTable, error) {
	s, err := service.NewSubTable(conf, opt...)
	if err != nil {
		return nil, err
	}
	return &SubTable{
		subTable: s,
	},nil
}
// CreateSubTable CreateSubTable
func (s *SubTable)CreateSubTable(c *gin.Context)  {
	req := &service.CreateSubTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(s.subTable.CreateSubTable(logger.CTXTransfer(c), req)).Context(c)
}

// UpdateSubTable UpdateSubTable
func (s *SubTable)UpdateSubTable(c *gin.Context)  {
	req := &service.UpdateSubTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(s.subTable.UpdateSubTable(logger.CTXTransfer(c), req)).Context(c)
}

// DeleteSubTable DeleteSubTable
func (s *SubTable)DeleteSubTable(c *gin.Context)  {
	req := &service.DeleteSubTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(s.subTable.DeleteSubTable(logger.CTXTransfer(c), req)).Context(c)
}

// GetSubTable GetSubTable
func (s *SubTable)GetSubTable(c *gin.Context)  {
	req := &service.GetSubTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(s.subTable.GetSubTableByTableID(logger.CTXTransfer(c), req)).Context(c)
}

// GetSubTables GetSubTables
func (s *SubTable)GetSubTables(c *gin.Context)  {
	req := &service.GetSubTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(s.subTable.GetSubTableByCondition(logger.CTXTransfer(c), req)).Context(c)
}

// GetSubTablesByType GetSubTables by type
func (s *SubTable)GetSubTablesByType(c *gin.Context)  {
	req := &service.GetSubTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(s.subTable.GetSubTableByType(logger.CTXTransfer(c), req)).Context(c)
}