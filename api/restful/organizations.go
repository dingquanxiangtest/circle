package restful

import (
	"net/http"

	"git.internal.yunify.com/qxp/misc/header2"
	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/resp"
	"git.internal.yunify.com/qxp/molecule/pkg/client"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/gin-gonic/gin"
)

// Organizations Organizations
type Organizations struct {
	Organizations client.Organizations
}

// NewOrganizations NewOrganizations
func NewOrganizations(c *config.Config) *Organizations {
	return &Organizations{
		Organizations: client.NewOrganizations(c.InternalNet),
	}
}

// DEPTree DEPTree
func (o *Organizations) DEPTree(c *gin.Context) {

	req := &client.DEPTreeReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(o.Organizations.DEPTree(logger.CTXTransfer(c), req)).Context(c)

}

// SelectDepByIDs SelectDepByIDs
func (o *Organizations) SelectDepByIDs(c *gin.Context) {
	req := &client.SelectDepByIDsReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(o.Organizations.SelectDepByIDs(logger.CTXTransfer(c), req)).Context(c)
}

// OnlineUserDep OnlineUserDep
func (o *Organizations) OnlineUserDep(c *gin.Context) {
	profile := header2.GetProfile(c)
	depID := profile.DepartmentID

	req := &client.SelectDepByIDsReq{
		IDs: []string{depID},
	}
	resp.Format(o.Organizations.SelectDepByIDs(logger.CTXTransfer(c), req)).Context(c)
}

// UserUsersInfo UserUsersInfo
func (o *Organizations) UserUsersInfo(c *gin.Context) {

	req := &client.UserSelectByIDsReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(o.Organizations.UserSelectByIDs(logger.CTXTransfer(c), req)).Context(c)

}

// OnlineUserInfo OnlineUserInfo
func (o *Organizations) OnlineUserInfo(c *gin.Context) {
	profile := header2.GetProfile(c)
	userID := profile.UserID

	req := &client.UserSelectByIDsReq{
		IDs: []string{userID},
	}
	resp.Format(o.Organizations.UserSelectByIDs(logger.CTXTransfer(c), req)).Context(c)
}

// SelectUserByCondition SelectUserByCondition
func (o *Organizations) SelectUserByCondition(c *gin.Context) {
	req := &client.AdminSelectByConditionReq{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(o.Organizations.AdminSelectByCondition(logger.CTXTransfer(c), req)).Context(c)
}
