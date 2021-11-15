package restful

import (
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"git.internal.yunify.com/qxp/molecule/internal/comet"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Process Process engine
type Process struct {
	process comet.Process
}

// NewProcess new NewProcess
func NewProcess(conf *config.Config, opt ...service.Options) (*Process, error) {
	p, err := comet.NewProcessServer(conf, opt...)
	if err != nil {
		return nil, err
	}
	return &Process{
		process: p,
	}, nil
}

// GetSubTable GetSubTable
func (p *Process) GetSubTable(c *gin.Context) {
	req := &service.GetSubTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(p.process.GetSubTable(logger.CTXTransfer(c), req)).Context(c)
}
// GetSchema GetSchema
func (p *Process) GetSchema(c *gin.Context) {
	req := &service.GetTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(p.process.GetTableByID(logger.CTXTransfer(c), req)).Context(c)
}

// GetRawSchema GetRawSchema
func (p *Process) GetRawSchema(c *gin.Context) {
	req := &comet.ProcessTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(p.process.GetRawTableByID(logger.CTXTransfer(c), req)).Context(c)
}

// GetData GetData
func (p *Process) GetData(c *gin.Context) {
	req := &comet.GetDataReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.process.GetData(logger.CTXTransfer(c), req)).Context(c)
}

// BatchGetData BatchGetData
func (p *Process) BatchGetData(c *gin.Context) {
	req := &comet.BatchGetDataReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(p.process.BatchGetData(logger.CTXTransfer(c), req)).Context(c)
}

// UpdateProcessData UpdateProcessData
func (p *Process) UpdateProcessData(c *gin.Context) {
	req := &comet.CreateProcessDataReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(p.process.UpdateProcessData(logger.CTXTransfer(c),c, req)).Context(c)
}

// CreateProcessData CreateProcessData
func (p *Process) CreateProcessData(c *gin.Context) {
	req := &comet.CreateProcessDataReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(p.process.CreateProcessData(logger.CTXTransfer(c),c, req)).Context(c)
}