package restful

import (
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// DataSet DataSet
type DataSet struct {
	dataset service.DataSet
}

// NewDataSet 初始化
func NewDataSet(conf *config.Config, opt ...service.Options) (*DataSet, error) {
	d, err := service.NewDataSet(conf, opt...)
	if err != nil {
		return nil, err
	}
	return &DataSet{
		dataset: d,
	}, nil
}

// CreateDataSet CreateDataSet
func (d *DataSet) CreateDataSet(c *gin.Context) {
	req := &service.CreateDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.CreateDataSet(logger.CTXTransfer(c), req)).Context(c)

}

// GetDataSet GetDataSet
func (d *DataSet) GetDataSet(c *gin.Context) {
	req := &service.GetDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.GetDataSet(logger.CTXTransfer(c), req)).Context(c)
}

// UpdateDataSet UpdateDataSet
func (d *DataSet) UpdateDataSet(c *gin.Context) {
	req := &service.UpdateDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.UpdateDataSet(logger.CTXTransfer(c), req)).Context(c)

}

// GetByConditionSet GetByConditionSet
func (d *DataSet) GetByConditionSet(c *gin.Context) {
	req := &service.GetByConditionSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.GetByConditionSet(logger.CTXTransfer(c), req)).Context(c)
}

// DeleteDataSet DeleteDataSet
func (d *DataSet) DeleteDataSet(c *gin.Context) {
	req := &service.DeleteDataSetReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(d.dataset.DeleteDataSet(logger.CTXTransfer(c), req)).Context(c)

}
