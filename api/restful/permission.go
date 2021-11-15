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

// Permission per
type Permission struct {
	permission service.Permission
	filter     service.Filter
	menu       service.Menu
	//appCenterClient  client.AppCenter
}

// NewPermission NewPermission
func NewPermission(conf *config.Config, opt ...service.Options) (*Permission, error) {
	permission, err := service.NewPermission(conf, opt...)
	if err != nil {
		return nil, err
	}
	filter, err := service.NewFilter(conf, opt...)
	if err != nil {
		return nil, err
	}
	menu, err := service.New(conf, opt...)
	if err != nil {
		return nil, err
	}
	return &Permission{
		permission: permission,
		filter:     filter,
		menu:       menu,
	}, nil
}

// CreatePermissionGroup 创建一个用户组
func (per *Permission) CreatePermissionGroup(c *gin.Context) {
	req := &service.CreatePermissionGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	req.AppID = c.Param("appID")
	if req.AppID == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	resp.Format(per.permission.CreatePermissionGroup(logger.CTXTransfer(c), req)).Context(c)
}

// UpdatePermissionGroup UpdatePermissionGroup
func (per *Permission) UpdatePermissionGroup(c *gin.Context) {
	req := &service.UpdatePermissionGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.UpdatePermissionGroup(logger.CTXTransfer(c), req)).Context(c)
}

// DeletePermissionGroup 删除
func (per *Permission) DeletePermissionGroup(c *gin.Context) {

	req := &service.DeletePermissionGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.DeletePermissionGroup(logger.CTXTransfer(c), req)).Context(c)
}

// GetByIDPermissionGroup get data by id
func (per *Permission) GetByIDPermissionGroup(c *gin.Context) {

	req := &service.GetByIDPermissionGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetByIDPermissionGroup(logger.CTXTransfer(c), req)).Context(c)
}

// GetListPermissionGroup GetListPermissionGroup
func (per *Permission) GetListPermissionGroup(c *gin.Context) {
	req := &service.GetListPermissionGroupReq{}
	req.AppID = c.Param("appID")
	if req.AppID == "" {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetListPermissionGroup(logger.CTXTransfer(c), req)).Context(c)
}

// GetByConditionPerGroup get perGroup by condition
func (per *Permission) GetByConditionPerGroup(c *gin.Context) {
	req := &service.GetByConditionPerGroupReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetByConditionPerGroup(logger.CTXTransfer(c), req)).Context(c)
}

// GetDataAccessPermission GetDataAccessPermission
func (per *Permission) GetDataAccessPermission(c *gin.Context) {
	req := &service.GetDataAccessPermissionReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetDataAccessPermission(logger.CTXTransfer(c), req)).Context(c)
}

// GetOperatePermission GetOperatePermission
func (per *Permission) GetOperatePermission(c *gin.Context) {
	req := &service.GetOperatePermissionReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetOperatePermission(logger.CTXTransfer(c), req)).Context(c)
}

// GetPerOption GetPerOption
func (per *Permission) GetPerOption(c *gin.Context) {
	userID := header2.GetProfile(c).UserID
	depID := header2.GetDepartments(c)[len(header2.GetDepartments(c))-1]
	appID := c.Param("appID")
	ctx := logger.CTXTransfer(c)
	selectPer, err := per.permission.GetPerSelect(ctx, &service.GetPerSelectReq{
		AppID:  appID,
		UserID: userID,
		DepID:  depID,
	})

	if err != nil {
		logger.Logger.Errorw(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	optionPer, err := per.permission.GetGroupPerByUserInfo(ctx, &service.GetGroupPerByUserInfoReq{
		AppID:  appID,
		DepID:  depID,
		UserID: userID,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(map[string]interface{}{
		"selectPer": selectPer,
		"optionPer": optionPer.PerGroupArr,
	}, nil).Context(c)
}

// VisibilityApp VisibilityApp
func (per *Permission) VisibilityApp(c *gin.Context) {
	req := &service.VisibilityAppReq{}

	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.VisibilityApp(logger.CTXTransfer(c), req)).Context(c)
}

// SaveForm SaveForm
func (per *Permission) SaveForm(c *gin.Context) {
	req := &service.SaveFormReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.SaveForm(logger.CTXTransfer(c), req)).Context(c)
}

// DeleteForm DeleteForm
func (per *Permission) DeleteForm(c *gin.Context) {
	req := &service.DeleteFormReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.DeleteForm(logger.CTXTransfer(c), req)).Context(c)

}

// GetForm GetForm
func (per *Permission) GetForm(c *gin.Context) {
	req := &service.GetFormReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetForm(logger.CTXTransfer(c), req)).Context(c)
}

// GetPerData GetPerData
func (per *Permission) GetPerData(c *gin.Context) {
	req := &struct {
		PerGroupID string `json:"perGroupID"`
		FormID     string `json:"formID"`
	}{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx := logger.CTXTransfer(c)
	dataAccessResp, err := per.permission.GetDataAccessPermission(ctx, &service.GetDataAccessPermissionReq{
		FormID:     req.FormID,
		PerGroupID: req.PerGroupID,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	operateResp, err := per.permission.GetOperatePermission(ctx, &service.GetOperatePermissionReq{
		FormID:     req.FormID,
		PerGroupID: req.PerGroupID,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	filterResp, err := per.filter.GetJSONFilter(ctx, &service.GetFilterReq{
		FormID:     req.FormID,
		PerGroupID: req.PerGroupID,
	})
	if err != nil {
		logger.Logger.Errorw(err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(map[string]interface{}{
		"opt":        operateResp,
		"dataAccess": dataAccessResp,
		"filter":     filterResp,
	}, nil).Context(c)

}

// GetOperate GetOperate
func (per *Permission) GetOperate(c *gin.Context) {
	req := &service.GetOperateReq{}
	req.UserID = header2.GetProfile(c).UserID
	req.DepID = header2.GetDepartments(c)[len(header2.GetDepartments(c))-1]
	req.AppID = c.Param("appID")
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetOperate(logger.CTXTransfer(c), req)).Context(c)

}

// SaveUserPerMatch SaveUserPerMatch
func (per *Permission) SaveUserPerMatch(c *gin.Context) {
	req := &service.SaveUserPerMatchReq{}
	req.UserID = header2.GetProfile(c).UserID
	req.AppID = c.Param("appID")
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.SaveUserPerMatch(logger.CTXTransfer(c), req)).Context(c)
}

// UpdatePerName UpdatePerName
func (per *Permission) UpdatePerName(c *gin.Context) {
	req := &service.UpdatePerNameReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.UpdatePerName(logger.CTXTransfer(c), req)).Context(c)
}

// AddAppDepUser AddAppDepUser
func (per *Permission) AddAppDepUser(c *gin.Context) {
	req := &service.VisibilityAppReq{}
	req.AppID = c.Param("appID")
	_, err := per.permission.VisibilityApp(logger.CTXTransfer(c), req)
	if err != nil {
		logger.Logger.Errorw(err.Error(), logger.STDRequestID(c))
	}
}

// SavePagePermission SavePagePermission
func (per *Permission) SavePagePermission(c *gin.Context) {
	req := &service.UpdatePagePermissionReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.UpdatePagePermission(logger.CTXTransfer(c), req)).Context(c)
}

// GetPagePermission GetPagePermission
func (per *Permission) GetPagePermission(c *gin.Context) {
	req := &service.GetGroupPageReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.GetGroupPage(logger.CTXTransfer(c), req)).Context(c)
}

// ModifyPagePermission modify the permission group page permission
func (per *Permission) ModifyPagePermission(c *gin.Context) {
	req := &service.ModifyPagePerReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(per.permission.ModifyPagePer(logger.CTXTransfer(c), req)).Context(c)
}

// GetPerGroupsByMenu Get the permission groups by menu id.
func (per *Permission) GetPerGroupsByMenu(c *gin.Context) {
	req := &service.GetPerGroupByMenuReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx := logger.CTXTransfer(c)
	m, err := per.menu.FindByID(ctx, &service.FindByIDReq{
		ID: req.MenuID,
	})

	if err != nil {
		resp.Format(nil, err).Context(c)
	}
	req.MenuType = m.MenuType
	req.AppID = c.Param("appID")
	resp.Format(per.permission.GetPerGroupByMenu(ctx, req)).Context(c)
}
