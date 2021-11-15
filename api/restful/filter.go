package restful

import (
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Filter filter
type Filter struct {
	filter service.Filter
}

// NewFilter 初始化
func NewFilter(conf *config.Config, opt ...service.Options) (*Filter, error) {
	c, err := service.NewFilter(conf, opt...)
	if err != nil {
		return nil, err
	}
	return &Filter{
		filter: c,
	}, nil
}

// Save 存储过滤字段
func (f *Filter) Save(c *gin.Context) {
	req := &service.SaveFilterReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(f.filter.SaveJSONFilter(logger.CTXTransfer(c), req)).Context(c)
}

// Get 获取过滤字段
func (f *Filter) Get(c *gin.Context) {
	req := &service.GetFilterReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(f.filter.GetJSONFilter(logger.CTXTransfer(c), req)).Context(c)
}
