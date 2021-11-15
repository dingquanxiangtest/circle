package service

import (
	"git.internal.yunify.com/qxp/molecule/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
)

// Opt options interface
type Opt interface {
	SetMongo(client *mongo.Client, dbName string)
}

// Options type options functions]
type Options func(Opt)

// WithMongo set mongo client to OPT
func WithMongo(client *mongo.Client, dbName string) Options {
	return func(o Opt) {
		o.SetMongo(client, dbName)
	}
}

// WithPermission WithPermission
func WithPermission(repo models.PermissionRepo) Options {
	return func(o Opt) {
		switch o := o.(type) {
		case *permission:
			o.permissionRepo = repo
		}

	}
}
