package main

import (
	"context"
	"encoding/json"
	"flag"
	"git.internal.yunify.com/qxp/misc/logger"
	mongo2 "git.internal.yunify.com/qxp/misc/mongo"
	filters2 "git.internal.yunify.com/qxp/molecule/internal/filters"
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"git.internal.yunify.com/qxp/molecule/internal/models/mongo"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	configPath = flag.String("config", "configs/config.yml", "-config 配置文件地址")
)

//处理旧版本的已经保存的过滤规则
func main() {
	flag.Parse()

	conf, err := config.NewConfig(*configPath)
	if err != nil {
		panic(err)
	}
	client, err := mongo2.New(&conf.Mongo)
	if err != nil {
		panic(err)
	}
	err = logger.New(&conf.Log)
	if err != nil {
		panic(err)
	}
	db := client.Database(conf.Service.DB)
	filters := make([]models.Filter, 0)
	ctx := context.Background()
	m := bson.M{}
	opts := &options.FindOptions{}
	cursor, err := db.Collection("filter").Find(ctx, m, opts)
	err = cursor.All(ctx, &filters)
	if err != nil {
		panic(err)
	}
	for k := range filters {
		if filters[k].FieldJSON != "" && filters[k].FieldJSON != "null" {
			schema := filters2.Schema{}
			err = json.Unmarshal([]byte(filters[k].WebSchema), &schema)
			if err != nil {
				logger.Logger.Errorw("to json fail", filters[k], err)
				continue
			}
			newFilter := filters2.DealSchemaToFilterType(schema)
			marshal, _ := json.Marshal(newFilter)
			filters[k].FieldJSON = string(marshal)
			repo := mongo.NewFilterRepo()
			err = repo.Update(ctx, db, &filters[k])
			if err != nil {
				logger.Logger.Errorw(filters[k].ID, filters[k], err)
				continue
			}
		}
	}
	logger.Logger.Info("处理数据：", len(filters), "条")
	logger.Logger.Info("job is done!")
}
