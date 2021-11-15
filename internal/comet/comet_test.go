package comet

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.internal.yunify.com/qxp/misc/logger"
	"git.internal.yunify.com/qxp/misc/mongo"

	"git.internal.yunify.com/qxp/molecule/internal/dorm"
	"git.internal.yunify.com/qxp/molecule/internal/filters"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CometSuite struct {
	suite.Suite

	ctx    context.Context
	userID string

	opt  service.Options
	conf *config.Config

	r *gin.Engine

	comet      Comet
	permission service.Permission
	filter     service.Filter
	FormID     string
	AppID      string
	misc       map[string]interface{}
}

func _TestComet(t *testing.T) {
	suite.Run(t, new(CometSuite))
}

func (suite *CometSuite) SetupTest() {
	suite.userID = "1"
	suite.FormID = "1"
	suite.AppID = "1"
	suite.misc = make(map[string]interface{})
	suite.ctx = logger.GenRequestID(context.TODO())

	var err error
	suite.conf, err = config.NewConfig("../../configs/config.yml")
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.conf)

	err = logger.New(&suite.conf.Log)
	assert.Nil(suite.T(), err)

	client, err := mongo.New(&suite.conf.Mongo)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), client)

	suite.opt = service.WithMongo(client, suite.conf.Service.DB)

}

func (suite *CometSuite) AfterTest(suiteName, testName string) {
}

func (suite *CometSuite) BeforeTest(suiteName, testName string) {
}

func (suite *CometSuite) PermissionBefore() {
	var err error
	suite.permission, err = service.NewPermission(suite.conf, suite.opt)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.permission)

	permission, err := suite.permission.CreatePermissionGroup(suite.ctx, &service.CreatePermissionGroupReq{
		Name: "test",

		AppID: suite.AppID,
	})
	assert.Nil(suite.T(), err)
	require.NotNil(suite.T(), permission)
	suite.misc["permissionID"] = permission.ID
	//
	//_, err = suite.permission.SaveOperatePermission(suite.ctx, &service.SaveOperatePermissionReq{
	//	PerGroupID: permission.ID,
	//	Authority:  models.OPCreate | models.OPUpdate | models.OPRead | models.OPDelete,
	//})
	assert.Nil(suite.T(), err)

	_, err = suite.permission.UpdatePermissionGroup(suite.ctx, &service.UpdatePermissionGroupReq{
		ID: permission.ID,
		Scopes: []*models.ScopesVO{{
			Type: 1,
			ID:   suite.userID,
			Name: "alice",
		}},
	})
	assert.Nil(suite.T(), err)

}

func (suite *CometSuite) PermissionAfter() {
	permissionID, ok := suite.misc["permissionID"].(string)
	assert.Equal(suite.T(), true, ok)
	assert.NotEqual(suite.T(), "", permissionID)

	_, err := suite.permission.DeletePermissionGroup(suite.ctx, &service.DeletePermissionGroupReq{
		ID: permissionID,
	})
	assert.Nil(suite.T(), err)
	//delFilterReq := service.DELFilterReq{
	//	PermissionGroupID: permissionID,
	//}
	//err = suite.filter.Delete(suite.ctx, &delFilterReq)
	assert.Nil(suite.T(), err)
}

func (suite *CometSuite) TestComet() {
	suite.PermissionBefore()
	defer suite.PermissionAfter()

	var err error
	suite.comet, err = New(suite.conf,
		service.WithPermission(
			models.NewPerMissionMock(
				models.WithPermissionMockID(suite.misc["permissionID"].(string)))),
		suite.opt)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.comet)

	suite.filter, err = service.NewFilter(suite.conf, suite.opt)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), suite.filter)

	suite.SaveFilter()

	suite.r = gin.Default()
	suite.r.POST("/api/v1/:tableName", suite.comet.Handle)

	var request = func(suite *CometSuite, body []byte) *httptest.ResponseRecorder {
		var uri = "/api/v1/1"

		req := httptest.NewRequest("POST", uri, bytes.NewBuffer(body))
		req.Header.Set("User-Id", suite.userID)
		w := httptest.NewRecorder()
		suite.r.ServeHTTP(w, req)

		return w
	}

	type entity struct {
		ID      string      `json:"_id"`
		Name    string      `json:"name"`
		Age     int         `json:"age"`
		Address interface{} `json:"address"`
	}

	createReq, err := json.Marshal(&input{
		Method: "create",
		Entity: entity{
			Name: "alice",
			Age:  18,
			Address: struct {
				Country  string `json:"country"`
				Province string `json:"province"`
				City     string `json:"city"`
			}{
				Country:  "china",
				Province: "sichuan",
				City:     "chengdu",
			},
		}})
	assert.Nil(suite.T(), err)

	w := request(suite, createReq)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	findReq, err := json.Marshal(map[string]interface{}{
		"method": "find",
		"condition": []dorm.Condition{{
			Key:   "name",
			Op:    "like",
			Value: []interface{}{"alice"},
		}},
		"page": 1,
		"size": 10,
		"sort": []string{"-name", "age"},
	})
	assert.Nil(suite.T(), err)

	w = request(suite, findReq)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	findBody, err := io.ReadAll(w.Body)
	assert.Nil(suite.T(), err)

	find := &struct {
		Code int `json:"code"`
		Data struct {
			Entities []entity `json:"entities"`
			Total    int
		} `json:"data"`
	}{}

	err = json.Unmarshal(findBody, &find)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), 0, find.Code)
	assert.Equal(suite.T(), 1, len(find.Data.Entities))

	entities := find.Data.Entities

	entityID := entities[0].ID
	assert.NotEqual(suite.T(), "", entityID)

	findOneReq, err := json.Marshal(map[string]interface{}{
		"method": "findOne",
		"condition": []dorm.Condition{{
			Key:   "_id",
			Op:    "like",
			Value: []interface{}{entityID},
		}},
	})
	assert.Nil(suite.T(), err)

	w = request(suite, findOneReq)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	updateSetReq, err := json.Marshal(map[string]interface{}{
		"method": "update",
		"condition": []dorm.Condition{{
			Key:   "_id",
			Op:    "like",
			Value: []interface{}{entityID},
		}},
		"entity": map[string]interface{}{
			"age": 19,
		},
	})
	assert.Nil(suite.T(), err)
	w = request(suite, updateSetReq)
	assert.Equal(suite.T(), http.StatusOK, w.Code)

	DeleteReq, err := json.Marshal(map[string]interface{}{
		"method": "delete",
		"condition": []dorm.Condition{{
			Key:   "_id",
			Op:    "like",
			Value: []interface{}{entityID},
		}},
	})
	assert.Nil(suite.T(), err)

	w = request(suite, DeleteReq)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *CometSuite) SaveFilter() {
	//提交字段吞吐
	schemaJSON := filters.Schema{
		Title: "人员",
		Types: "object",
		XInternal: filters.XInternal{
			Sortable:   true,
			Permission: 2,
		},
	}
	properties := make(map[string]filters.Schema)

	properties["_id"] = filters.Schema{
		Types: "string",
		XInternal: filters.XInternal{
			Sortable:   true,
			Permission: 2,
		},
	}
	properties["name"] = filters.Schema{
		Types: "string",
		XInternal: filters.XInternal{
			Sortable:   true,
			Permission: 2,
		},
	}
	address := make(map[string]filters.Schema)
	address["_id"] = filters.Schema{
		Types: "string",
		XInternal: filters.XInternal{
			Sortable:   true,
			Permission: 2,
		},
	}
	address["country"] = filters.Schema{
		Types: "string",
		XInternal: filters.XInternal{
			Sortable:   true,
			Permission: 2,
		},
	}
	address["province"] = filters.Schema{
		Types: "string",
		XInternal: filters.XInternal{
			Sortable:   true,
			Permission: 2,
		},
	}
	address["city"] = filters.Schema{
		Types: "string",
		XInternal: filters.XInternal{
			Sortable:   true,
			Permission: 2,
		},
	}
	properties["age"] = filters.Schema{
		Types: "number",
		XInternal: filters.XInternal{
			Sortable:   true,
			Permission: 2,
		},
	}
	properties["address"] = filters.Schema{
		Types: "object",
		XInternal: filters.XInternal{
			Sortable:   true,
			Permission: 2,
		},
		Properties: address,
	}
	schemaJSON.Properties = properties

	filterReq := service.SaveFilterReq{
		//PermissionGroupID: suite.misc["permissionID"].(string),
		Schema: schemaJSON,
	}
	filterResp, err := suite.filter.SaveJSONFilter(suite.ctx, &filterReq)
	assert.Nil(suite.T(), err)
	assert.Nil(suite.T(), filterResp)
}
