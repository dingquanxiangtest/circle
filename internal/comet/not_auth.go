package comet

import (
	"git.internal.yunify.com/qxp/molecule/internal/dorm"
	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"
	"git.internal.yunify.com/qxp/molecule/internal/service"
	"git.internal.yunify.com/qxp/molecule/pkg/misc/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type noAuth struct {
	comet1
}

// NewNoAuth NewNoAuth
func NewNoAuth(conf *config.Config, opts ...service.Options) (Plugs, error) {
	p, err := service.NewPermission(conf, opts...)
	if err != nil {
		return nil, err
	}
	filter, err := service.NewFilter(conf, opts...)
	if err != nil {
		return nil, err
	}
	t, err := service.NewKernel(conf, opts...)
	if err != nil {
		return nil, err
	}
	a := &noAuth{
		comet1: comet1{
			cm: &CMongo{
				dc:    clause.New(),
				query: dorm.NewQuery(),
				ag:    clause.NewAg(),
			},
			permission: p,
			filter:     filter,
			schema:     t,
		},
	}
	for _, opt := range opts {
		opt(a)
	}
	return a, nil
}

func (c1 *noAuth) SetMongo(client *mongo.Client, dbName string) {
	c1.cm.DB = client.Database(dbName)
}
