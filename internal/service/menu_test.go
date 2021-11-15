package service

import (
	"context"
	"testing"

	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/mongo"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MenuSuite struct {
	suite.Suite

	ctx context.Context

	opt  Options
	conf *config.Config
	menu Menu

	AppID string
	misc  map[string]interface{}
}

func TestMenuSuite(t *testing.T) {
	suite.Run(t, new(MenuSuite))
}

func (suite *MenuSuite) SetupTest() {
	suite.AppID = "test001"
	suite.misc = make(map[string]interface{})

	suite.ctx = logger.GenRequestID(context.TODO())
	var err error
	suite.conf, err = config.NewConfig("../../configs/config.yml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.conf)

	client, err := mongo.New(&suite.conf.Mongo)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), client)

	suite.opt = WithMongo(client, suite.conf.Service.DB)
	assert.NotNil(suite.T(), suite.opt)

}

func (suite *MenuSuite) MenuBefore() {

	var err error
	suite.menu, err = New(suite.conf, suite.opt)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.menu)

	var group *CreateGroupResp
	group, err = suite.menu.CreateGroup(suite.ctx, &CreateGroupReq{
		AppID:   suite.AppID,
		Name:    "分组测试",
		GroupID: "",
	})

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), group)
	suite.misc["groupID"] = group.ID

	var menu *CreateMenuResp
	menu, err = suite.menu.CreateMenu(suite.ctx, &CreateMenuReq{
		AppID:    suite.AppID,
		Name:     "页面测试",
		Icon:     "icon",
		Describe: "describe",
		GroupID:  "",
	})

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), menu)
	suite.misc["menuID"] = menu.ID

}

func (suite *MenuSuite) MenuAfter() {

	groupID, ok := suite.misc["groupID"].(string)
	assert.Equal(suite.T(), true, ok)
	assert.NotEqual(suite.T(), "", groupID)

	menuID, ok := suite.misc["menuID"].(string)
	assert.Equal(suite.T(), true, ok)
	assert.NotEqual(suite.T(), "", menuID)

	_, err := suite.menu.DeleteMenu(suite.ctx, &DeleteMenuReq{
		ID:      menuID,
		Sort:    1,
		GroupID: groupID,
		AppID:   suite.AppID,
	})
	assert.Nil(suite.T(), err)

	_, err = suite.menu.DeleteGroup(suite.ctx, &DeleteGroupReq{
		AppID:   suite.AppID,
		ID:      groupID,
		Sort:    1,
		GroupID: "",
	})
	assert.Nil(suite.T(), err)
}

func (suite *MenuSuite) TestMenu() {

	suite.MenuBefore()
	defer suite.MenuAfter()
	var err error

	menuID, ok := suite.misc["menuID"].(string)
	assert.Equal(suite.T(), true, ok)
	assert.NotEqual(suite.T(), "", menuID)
	var newMenuName = "new menu name"
	_, err = suite.menu.UpdateMenu(suite.ctx, &UpdateMenuReq{
		ID:       menuID,
		Name:     newMenuName,
		Icon:     "new icon",
		Describe: "new describe",
		AppID:    suite.AppID,
		GroupID:  "",
	})
	assert.Nil(suite.T(), err)
	suite.misc["newMenuName"] = newMenuName

	groupID, ok := suite.misc["groupID"].(string)
	assert.Equal(suite.T(), true, ok)
	assert.NotEqual(suite.T(), "", groupID)
	_, err = suite.menu.UpdateMenu(suite.ctx, &UpdateMenuReq{
		ID:       groupID,
		Name:     "new group name",
		Icon:     "",
		Describe: "",
		AppID:    suite.AppID,
		GroupID:  "",
	})
	assert.Nil(suite.T(), err)

	groups, err := suite.menu.ListAllGroup(suite.ctx, &ListAllGroupReq{
		AppID: suite.AppID,
	})
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), groups)

	elements, err := suite.menu.ListAll(suite.ctx, &ListAllReq{
		AppID: suite.AppID,
	})
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), elements)

	newMenuName, ok = suite.misc["newMenuName"].(string)
	assert.Equal(suite.T(), true, ok)
	assert.NotEqual(suite.T(), "", newMenuName)

	_, err = suite.menu.Transfer(suite.ctx, &TransferReq{
		ID:          menuID,
		AppID:       suite.AppID,
		Name:        newMenuName,
		FromSort:    2,
		ToSort:      1,
		FromGroupID: "",
		ToGroupID:   "",
	})
	assert.Nil(suite.T(), err)

	_, err = suite.menu.Transfer(suite.ctx, &TransferReq{
		ID:          menuID,
		AppID:       suite.AppID,
		Name:        newMenuName,
		FromSort:    1,
		ToSort:      2,
		FromGroupID: "",
		ToGroupID:   "",
	})
	assert.Nil(suite.T(), err)

	_, err = suite.menu.Transfer(suite.ctx, &TransferReq{
		ID:          menuID,
		AppID:       suite.AppID,
		Name:        newMenuName,
		FromSort:    2,
		ToSort:      1,
		FromGroupID: "",
		ToGroupID:   groupID,
	})
	assert.Nil(suite.T(), err)

	//userMenu, err := suite.menu.UserListAll(suite.ctx, &UserListAllReq{
	//	AppID:   suite.AppID,
	//	//MenusID: []string{menuID},
	//})
	//assert.NotNil(suite.T(), userMenu)
	//assert.Nil(suite.T(), err)
}
