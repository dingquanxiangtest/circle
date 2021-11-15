package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"git.internal.yunify.com/qxp/misc/client"
	"git.internal.yunify.com/qxp/misc/id2"
	"git.internal.yunify.com/qxp/misc/logger"
	mongo2 "git.internal.yunify.com/qxp/misc/mongo"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"git.internal.yunify.com/qxp/molecule/internal/service/swagger"
	client2 "git.internal.yunify.com/qxp/molecule/pkg/client"
	"net/http"

	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/olekukonko/tablewriter"
	"os"
)

var (
	configPath = flag.String("config", "configs/config.yml", "-config 配置文件地址")
)

const (
	polyapiHost = "http://polyapi:9090/api/v1/polyapi/inner/requestPoly"

	parentName = "/system/poly/permissionInit"
)

// schema
func main() {
	flag.Parse()

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}
	client1, err := mongo2.New(&conf.Mongo)
	if err != nil {
		panic(err)
	}

	err = logger.New(&conf.Log)
	if err != nil {
		panic(err)
	}
	db := client1.Database(conf.Service.DB)
	tables := make([]models.Table, 0)
	ctx := context.Background()

	postClient := client.New(conf.InternalNet)
	// 查询schema
	filter := bson.M{}
	cursor, err := db.Collection("table_schema").Find(ctx, filter)
	err = cursor.All(ctx, &tables)
	sw := swagger.NewSW()
	apiClient := client2.NewPolyAPI(conf.InternalNet)
	if err != nil {
		panic(err)
	}

	var total = 0
	var success = 0
	var fail = 0
	var noDo = 0

	t1 := time.Now()
	arr := make([][]string, 0)
	for _, value := range tables {
		total = total + 1 // 总数加一
		// 根据tableID 查询
		database := &models.DataBaseSchema{}
		f := bson.M{
			"table_id": value.TableID,
		}
		err := db.Collection("database_schema").FindOne(ctx, f).Decode(database)
		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			// 去查询menu
			mf := bson.M{
				"_id": value.TableID,
			}
			t2 := time.Now()

			menu := &models.Menu{}
			err := db.Collection("menu").FindOne(ctx, mf).Decode(menu)
			menuTime := time.Since(t2)
			logger.Logger.Info("menu spend time", menuTime)
			if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
				noDo = noDo + 1
				continue
			}
			if err != nil {
				var str = []string{value.TableID, "get menu is fail ", err.Error()}
				arr = append(arr, str)
				fail = fail + 1
				continue
			}

			var s = map[string]interface{}{}
			if _, ok := value.Schema["properties"].(map[string]interface{}); !ok {

			}
			s = value.Schema["properties"].(map[string]interface{})

			t3 := time.Now()
			convert, total, err := swagger.Convert1(s)
			convertTime := time.Since(t3)
			logger.Logger.Info("convert spend time", convertTime)

			if err != nil {
				var str = []string{value.TableID, "convert  schema is fail ", err.Error()}
				arr = append(arr, str)
				fail = fail + 1
				continue
			}
			tableSchema := &models.DataBaseSchema{
				ID:          id2.GenID(),
				Title:       menu.Name,
				AppID:       menu.AppID,
				TableID:     value.TableID,
				FieldLen:    total,
				Description: menu.Describe,
				Source:      models.FormSource,
				Schema:      value.Schema,
			}
			t4 := time.Now()
			_, err = db.Collection("database_schema").InsertOne(ctx, &tableSchema)
			InsertTime := time.Since(t4)
			logger.Logger.Info("insert spend time", InsertTime)
			duration := time.Since(t2)
			logger.Logger.Info("get menu , convert , and  insert database_schema spend time:", duration)
			if err != nil {
				var str = []string{value.TableID, "insert database_schema    is fail ", err.Error()}
				arr = append(arr, str)
				fail = fail + 1
				logger.Logger.Errorw("tableID is init fail ", value.TableID, err.Error())
				continue
			}
			t1 := time.Now()
			// 先调用初始化权限
			perPoly(ctx, menu.AppID, "初始化权限组", "初始化权限组", 1, postClient)

			// 注册到yapi
			err = gen(ctx, convert, menu.AppID, value.TableID, apiClient, "form", sw)

			if err != nil {
				var str = []string{value.TableID, "register  api  is fail ", err.Error()}
				arr = append(arr, str)
				fail = fail + 1
				continue
			}
			since := time.Since(t1)

			logger.Logger.Info("register send time:", since)
			success = success + 1
			logger.Logger.Info("tableID is end", value.TableID)
			continue

		}
		if err != nil {
			var str = []string{value.TableID, "get database_schema is fail ", err.Error()}
			arr = append(arr, str)
			fail = fail + 1
			continue
		}
		var s = map[string]interface{}{}
		// 等于表单驱动的，只要去注册yapi ，不要重新更新
		if _, ok := value.Schema["properties"].(map[string]interface{}); !ok {

		}
		s = value.Schema["properties"].(map[string]interface{})
		convert, _, err := swagger.Convert1(s)
		if err != nil {
			var str = []string{value.TableID, "get Convert1 is fail ", err.Error()}
			arr = append(arr, str)
			fail = fail + 1
			continue
		}
		now := time.Now()

		perPoly(ctx, database.AppID, "初始化权限组", "初始化权限组", 1, postClient)
		content := "form"
		if database.Source == 2 {
			content = "custom"
		}
		err = gen(ctx, convert, database.AppID, value.TableID, apiClient, content, sw)
		if err != nil {
			var str = []string{value.TableID, "register  api  is fail ", err.Error()}
			arr = append(arr, str)
			fail = fail + 1
			continue
		}
		since := time.Since(now)
		fmt.Sprintf("%s", since)
		logger.Logger.Info(since)
		success = success + 1
	}
	since := time.Since(t1)
	logger.Logger.Info("time", since)

	Write([]string{"total", "success", "fail", "noDo"}, [][]string{{fmt.Sprintf("%d", total), fmt.Sprintf("%d", success), fmt.Sprintf("%d", fail), fmt.Sprintf("%d", noDo)}})
	Write([]string{"id", "errMsg", "err"}, arr)

}

func perPoly(ctx context.Context, appID, name, description string, types int64, postClient http.Client) error {
	params := struct {
		AppID       string `json:"appID"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Types       int64  `json:"types"`
	}{
		AppID:       appID,
		Name:        name,
		Description: description,
		Types:       types,
	}
	resp := struct {
	}{}
	err := client.POST(ctx, &postClient, polyapiHost+parentName, params, resp)
	if err != nil {
		return err
	}
	return nil
}

func gen(ctx context.Context, convert map[string]interface{}, appID, tableID string, apiClient client2.PolyAPI, content string, sw *swagger.Swagger) error {
	genSwagger, err := swagger.GenSwagger(convert, appID, tableID, sw)
	if err != nil {
		return err
	}
	swagger, err := json.Marshal(genSwagger)
	if err != nil {
		return err
	}
	_, err = apiClient.RegSwagger(ctx, "structor", string(swagger), appID, content)
	if err != nil {
		return err
	}
	return nil
}

// Write Write
func Write(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // Add Bulk Data
	table.Render()
}
