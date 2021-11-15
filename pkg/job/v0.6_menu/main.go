package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"git.internal.yunify.com/qxp/misc/logger"
	mongo2 "git.internal.yunify.com/qxp/misc/mongo"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"github.com/olekukonko/tablewriter"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	configPath = flag.String("config", "../../../configs/config.yml", "-config 配置文件地址")
)

// schema
func main() {
	flag.Parse()

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}
	mongoClient, err := mongo2.New(&conf.Mongo)
	if err != nil {
		panic(err)
	}

	err = logger.New(&conf.Log)
	if err != nil {
		panic(err)
	}
	db := mongoClient.Database(conf.Service.DB)

	menus := make([]*models.Menu, 0)
	ctx := context.Background()

	// 查询所有菜单
	startTime := time.Now()
	filter := bson.M{}
	cursor, err := db.Collection("menu").Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info("find all menus time consuming ", time.Since(startTime).Seconds())
	err = cursor.All(ctx, &menus)
	if err != nil {
		panic(err)
	}

	var (
		total   int
		success int
		fail    int
		noDo    int
	)

	errSet := make([][]string, 0)
	for _, value := range menus {
		total = total + 1 // 总数加一
		// 处理分组与自定义页面
		if value.MenuType != models.MenuType {
			if _, err := db.Collection("menu").UpdateByID(ctx, value.ID, bson.M{
				"$set": bson.M{
					"binding_state": models.Bound,
				},
			}); err != nil {
				fail++
				var res = []string{value.ID, "", err.Error()}
				errSet = append(errSet, res)
				continue
			}
			success++
			continue
		}

		// 根据menuID查询
		table := &models.Table{}
		err := db.Collection("table_schema").FindOne(ctx, bson.M{
			"table_id": value.ID,
		}).Decode(table)

		if err == mongo.ErrNilDocument || err == mongo.ErrNoDocuments {
			if _, err := db.Collection("menu").UpdateByID(ctx, value.ID, bson.M{
				"$set": bson.M{
					"binding_state": models.Unbound,
				},
			}); err != nil {
				var res = []string{value.ID, "update table binding state failed", err.Error()}
				errSet = append(errSet, res)
				fail++
				continue
			}

			success++
			continue
		}

		if err != nil {
			var res = []string{value.ID, "find table schema failed", err.Error()}
			errSet = append(errSet, res)
			fail++
			continue
		}

		if _, err = db.Collection("menu").UpdateByID(ctx, value.ID, bson.M{
			"$set": bson.M{
				"binding_state": models.Bound,
			},
		}); err != nil {
			var res = []string{value.ID, "update table binding state failed", err.Error()}
			errSet = append(errSet, res)
			fail++
			continue
		}
		success++
	}

	logger.Logger.Info("total time spent ", time.Since(startTime).Seconds())
	Write([]string{"total", "success", "fail", "noDo"}, [][]string{{fmt.Sprintf("%d", total), fmt.Sprintf("%d", success), fmt.Sprintf("%d", fail), fmt.Sprintf("%d", noDo)}})
	Write([]string{"id", "errMsg", "err"}, errSet)
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
