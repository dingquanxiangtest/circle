package restful

import (
	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CustomPage CustomPage
type CustomPage struct {
	customPage service.CustomPage
	menu       service.Menu
}

// NewCustomPage new a customPage manager
func NewCustomPage(conf *config.Config, opt ...service.Options) (*CustomPage, error) {
	c, err := service.NewCustomPage(conf, opt...)

	if err != nil {
		return nil, err
	}
	m, err := service.New(conf, opt...)
	if err != nil {
		return nil, err
	}
	return &CustomPage{
		customPage: c,
		menu:       m,
	}, nil
}

// CreateCustomPage CreateCustomPage
func (cus *CustomPage) CreateCustomPage(c *gin.Context) {
	req := &service.CreateCustomPageReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	req.AppID = c.Param("appID")
	profile := header2.GetProfile(c)
	req.UserID = profile.UserID
	req.UserName = profile.UserName

	cusResp, err := cus.customPage.CreateCustom(logger.CTXTransfer(c), req)
	if err != nil {
		resp.Format(nil, err)
	}
	_, err = cus.menu.ModifyMenuType(logger.CTXTransfer(c), &service.ModifyMenuTypeReq{ID: req.MenuID})
	resp.Format(cusResp, err).Context(c)
}

// UpdateCustomPage DeleteCustomPage
func (cus *CustomPage) UpdateCustomPage(c *gin.Context) {
	req := &service.UpdateCustomPageReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	profile := header2.GetProfile(c)
	req.UserID = profile.UserID
	req.UserName = profile.UserName
	resp.Format(cus.customPage.UpdateCustomPage(logger.CTXTransfer(c), req)).Context(c)
}

// GetByMenuID GetByMenuID
func (cus *CustomPage) GetByMenuID(c *gin.Context) {
	req := &service.GetByMenuIDReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(cus.customPage.GetByMenuID(logger.CTXTransfer(c), req)).Context(c)
}
