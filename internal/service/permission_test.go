package service

import (
	"context"
	"testing"

	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/mongo"

	"git.internal.yunify.com/qxp/molecule/internal/models"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PermissionSuite struct {
	suite.Suite

	ctx context.Context

	opt Options

	conf *config.Config

	r *gin.Engine

	permission Permission

	GroupID string

	FormID string

	AppID string

	UserID string
}

func _TestPermission(t *testing.T) {
	suite.Run(t, new(PermissionSuite))
}
func (suite *PermissionSuite) SetupSuite() {
	suite.ctx = logger.GenRequestID(context.TODO())
	suite.FormID = "1"
	suite.AppID = "1"
	suite.UserID = "1"
	var err error
	suite.conf, err = config.NewConfig("../../configs/config.yml")
	assert.Nil(suite.T(), err)           // 断言err 这个变量 ，会为空
	assert.NotNil(suite.T(), suite.conf) // 断言 conf 这个变量，不会为空
	err = logger.New(&suite.conf.Log)
	assert.Nil(suite.T(), err)
	client, err := mongo.New(&suite.conf.Mongo)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), client)
	suite.opt = WithMongo(client, suite.conf.Service.DB)
	suite.permission, err = NewPermission(suite.conf, suite.opt)

	pm := suite.permission.(*permission)
	pm.permissionRepo = models.NewPerMissionMock()
	suite.permission = pm

	assert.Nil(suite.T(), err) // 判断err 这个变量为空

}

func (suite *PermissionSuite) TestPermission() {
	// 创建权限用户组
	createPermissionGroupReq := &CreatePermissionGroupReq{

		AppID:       suite.AppID,
		Name:        "test",
		Description: "这是测试用例的数据",
	}
	resp, err := suite.permission.CreatePermissionGroup(suite.ctx, createPermissionGroupReq)
	assert.Nil(suite.T(), err) // 判断err 为空
	suite.GroupID = resp.ID
	// 修改权限用户组 名字和 描述
	updatePermissionGroupReq := &UpdatePermissionGroupReq{
		ID: suite.GroupID,
	}
	_, err = suite.permission.UpdatePermissionGroup(suite.ctx, updatePermissionGroupReq)
	assert.Nil(suite.T(), err)
	// 	给权限用户组添加人或者部门
	scopes := make([]*models.ScopesVO, 0)
	vo := &models.ScopesVO{
		ID:   suite.UserID,
		Type: 1,
		Name: "测试的人员",
	}
	scopes = append(scopes, vo)
	updatePermissionGroupReq = &UpdatePermissionGroupReq{
		ID:     suite.GroupID,
		Scopes: scopes,
	}
	_, err = suite.permission.UpdatePermissionGroup(suite.ctx, updatePermissionGroupReq)
	assert.Nil(suite.T(), err)
	//  根据id 得到权限用户组
	getByIDPermissionGroupReq := &GetByIDPermissionGroupReq{
		ID: suite.GroupID,
	}
	permission, err := suite.permission.GetByIDPermissionGroup(suite.ctx, getByIDPermissionGroupReq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), permission)
	//  得到权限列表
	getListPermissionGroupReq := &GetListPermissionGroupReq{}
	permissions, err := suite.permission.GetListPermissionGroup(suite.ctx, getListPermissionGroupReq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), permissions)
	// 根据条件权限用户组
	getConditionReq := &GetByConditionPerGroupReq{
		FormID: suite.FormID,
		UserID: suite.UserID,
	}
	per, err := suite.permission.GetByConditionPerGroup(suite.ctx, getConditionReq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), per)
	// 给appID 得到formID 数组
	getFormsPerGroupReq := &GetFormsPerGroupReq{
		//AppID: suite.AppID,
	}
	forms, err := suite.permission.GetFormsPerGroup(suite.ctx, getFormsPerGroupReq)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), forms)
	// 保存权限用户组，
	// 删除权限用户组
	deletePermissionGroupReq := &DeletePermissionGroupReq{
		ID: suite.GroupID,
	}
	_, err = suite.permission.DeletePermissionGroup(suite.ctx, deletePermissionGroupReq)
	assert.Nil(suite.T(), err)

}
