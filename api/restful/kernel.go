package restful

import (
	"net/http"

	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"

	"github.com/gin-gonic/gin"
)

// Kernel gin kernel
type Kernel struct {
	kernel service.Kernel
}

// NewKernel new kernel gin
func NewKernel(conf *config.Config, opt ...service.Options) (*Kernel, error) {
	k, err := service.NewKernel(conf, opt...)
	if err != nil {
		return nil, err
	}
	return &Kernel{
		kernel: k,
	}, nil
}

// CreateSchema 创建表
func (k *Kernel) CreateSchema(c *gin.Context) {
	req := &service.CreateTableReq{}
	req.AppID = c.Param("appID")
	req.UserName = header2.GetProfile(c).UserName
	req.UserID = header2.GetProfile(c).UserID
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	opts := []service.CreateSchemaOption{
		service.ConvertSchema(k.kernel),
		service.GenSwag(k.kernel),
	}
	resp.Format(k.kernel.CreateSchema(logger.CTXTransfer(c), req, opts...)).Context(c)
}

// CreateBlankSchema 创建表
func (k *Kernel) CreateBlankSchema(c *gin.Context) {
	req := &service.CreateBlankTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(k.kernel.CreateBlankSchema(logger.CTXTransfer(c), req)).Context(c)
}

// GetSchema GetSchema
func (k *Kernel) GetSchema(c *gin.Context) {
	req := &service.GetTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(k.kernel.GetSchemaByTableID(logger.CTXTransfer(c), req)).Context(c)
}

// DeleteSchema DeleteSchema
func (k *Kernel) DeleteSchema(c *gin.Context) {
	req := &service.DeleteTableReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(k.kernel.DeleteSchema(logger.CTXTransfer(c), req)).Context(c)
}

// CreateOrUpdateConfig table page config create or update
func (k *Kernel) CreateOrUpdateConfig(c *gin.Context) {
	req := &service.CreateConfigReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(k.kernel.CreateConfig(logger.CTXTransfer(c), req)).Context(c)
}

// DeleteConfig table page config delete
func (k *Kernel) DeleteConfig(c *gin.Context) {
	req := &service.DeleteConfigReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(k.kernel.DeleteConfig(logger.CTXTransfer(c), req)).Context(c)
}

// SearchSchema SearchSchema
func (k *Kernel) SearchSchema(c *gin.Context) {
	req := &service.SearchSchemaReq{}
	req.AppID = c.Param("appID")
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(k.kernel.SearchSchema(logger.CTXTransfer(c), req)).Context(c)
}

// GetXName GetXName
func (k *Kernel) GetXName(c *gin.Context) {
	req := &service.GetXNameReq{}
	req.AppID = c.Param("appID")
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(k.kernel.GetXName(logger.CTXTransfer(c), req)).Context(c)

}

// CheckRepeat CheckRepeat
func (k *Kernel) CheckRepeat(c *gin.Context) {
	req := &service.CheckRepeatReq{}
	req.AppID = c.Param("appID")
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(k.kernel.CheckRepeat(logger.CTXTransfer(c), req)).Context(c)

}

// GetModelDataByMenu GetModelDataByMenu
func (k *Kernel) GetModelDataByMenu(c *gin.Context) {
	req := &service.GetModelDataByMenuReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(k.kernel.GetModelDataByMenu(logger.CTXTransfer(c), req)).Context(c)
}
