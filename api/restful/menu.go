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

// Menu 菜单
type Menu struct {
	menu       service.Menu
	permission service.Permission
	customPage service.CustomPage
}

// NewMenu new a manager menu
func NewMenu(conf *config.Config, opt ...service.Options) (*Menu, error) {
	m, err := service.New(conf, opt...)
	if err != nil {
		return nil, err
	}
	p, err := service.NewPermission(conf, opt...)
	if err != nil {
		return nil, err
	}
	c, err := service.NewCustomPage(conf, opt...)
	if err != nil {
		return nil, err
	}
	return &Menu{
		menu:       m,
		permission: p,
		customPage: c,
	}, nil
}

// CreateMenu CreateMenu
func (m *Menu) CreateMenu(c *gin.Context) {
	req := &service.CreateMenuReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.menu.CreateMenu(logger.CTXTransfer(c), req)).Context(c)
}

// DeleteMenu DeleteMenu
func (m *Menu) DeleteMenu(c *gin.Context) {
	req := &service.DeleteMenuReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	_, err := m.menu.DeleteMenu(logger.CTXTransfer(c), req)
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	appID := c.Param("appID")
	// delete custom page
	cReq := &service.DeletePageMenuByMenuIDReq{
		MenuID: req.ID,
		AppID:  appID,
	}
	resp.Format(m.customPage.DeletePageMenuByMenuID(logger.CTXTransfer(c), cReq)).Context(c)
}

// UpdateMenu UpdateMenu
func (m *Menu) UpdateMenu(c *gin.Context) {
	req := &service.UpdateMenuReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.menu.UpdateMenu(logger.CTXTransfer(c), req)).Context(c)
}

//CreateGroup CreateGroup
func (m *Menu) CreateGroup(c *gin.Context) {
	req := &service.CreateGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.menu.CreateGroup(logger.CTXTransfer(c), req)).Context(c)
}

// DeleteGroup DeleteGroup
func (m *Menu) DeleteGroup(c *gin.Context) {
	req := &service.DeleteGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.menu.DeleteGroup(logger.CTXTransfer(c), req)).Context(c)
}

// ListAllGroup list all group
func (m *Menu) ListAllGroup(c *gin.Context) {
	req := &service.ListAllGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.menu.ListAllGroup(logger.CTXTransfer(c), req)).Context(c)
}

// ListAll list all menu
func (m *Menu) ListAll(c *gin.Context) {
	req := &service.ListAllReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.menu.ListAll(logger.CTXTransfer(c), req)).Context(c)
}

// Transfer Transfer
func (m *Menu) Transfer(c *gin.Context) {
	req := &service.TransferReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.menu.Transfer(logger.CTXTransfer(c), req)).Context(c)
}

// UserListAll UserListAll
func (m *Menu) UserListAll(c *gin.Context) {
	var (
		ctx    = logger.CTXTransfer(c)
		appID  = c.Param("appID")
		userID = header2.GetProfile(c).UserID
		depID  = header2.GetDepartments(c)[len(header2.GetDepartments(c))-1]
	)
	getFormResp, err := m.permission.GetFormsPerGroup(ctx, &service.GetFormsPerGroupReq{
		UserID: userID,
		AppID:  appID,
		DepID:  depID,
	})
	if err != nil {
		resp.Format(nil, err).Context(c)
		return
	}
	if getFormResp == nil {
		resp.Format(nil, nil).Context(c, http.StatusForbidden)
		return
	}
	resp.Format(m.menu.UserListAll(logger.CTXTransfer(c), &service.UserListAllReq{
		AppID:   appID,
		FormID:  getFormResp.FormID,
		PerType: getFormResp.PerType,
	})).Context(c)
}

// ListPage ListPage
func (m *Menu) ListPage(c *gin.Context) {
	req := &service.ListPageReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(m.menu.ListPage(logger.CTXTransfer(c), req)).Context(c)
}
