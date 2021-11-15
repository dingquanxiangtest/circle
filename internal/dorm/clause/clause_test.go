package clause

import (
	"context"
	"testing"

	"git.internal.yunify.com/qxp/misc/logger"
	mongo2 "git.internal.yunify.com/qxp/misc/mongo"
	"git.internal.yunify.com/qxp/misc/mysql2"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func newMONGO() (*mongo.Client, error) {
	return mongo2.New(&mongo2.Config{
		Direct: true,
		Hosts:  []string{"192.168.200.19:27017"},
		Credential: struct {
			AuthMechanism           string
			AuthMechanismProperties map[string]string
			AuthSource              string
			Username                string
			Password                string
			PasswordSet             bool
		}{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    "admin",
			Username:      "root",
			Password:      "uyWxtvt6gCOy3VPLB3rTpa0rQ",
			PasswordSet:   false,
		},
	})

}

func newMYSQL() (*gorm.DB, error) {
	err := logger.New(&logger.Config{
		Level:       -1,
		Development: false,
		Sampling: logger.Sampling{
			Initial:    100,
			Thereafter: 100,
		},
		OutputPath:      []string{"stderr"},
		ErrorOutputPath: []string{"stderr"},
	})
	if err != nil {
		return nil, err
	}

	db, err := mysql2.New(mysql2.Config{
		Host:     "192.168.200.18:3306",
		DB:       "test",
		User:     "root",
		Password: "uyWxtvt6gCOy3VPLB3rTpa0rQ",
		Log:      true,
	}, logger.Logger)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestMONGO(t *testing.T) {
	var (
		database   = "test"
		collection = "test"
		ctx        = context.Background()
	)

	client, err := newMONGO()
	if err != nil {
		t.Fatal(err)
		return
	}
	defer client.Disconnect(ctx)

	builder := &MONGO{}

	var expr Expressions = AND{}
	expr = expr.Set(
		"",
		IN{Column: "id", Values: []interface{}{"1", "2"}},
		LIKE{Column: "name", Values: "alice"},
	)

	expr.MongoBuild(builder)

	c := client.Database(database).Collection(collection)
	result := c.FindOne(ctx, builder.Vars)
	if result.Err() == mongo.ErrNoDocuments {
		err = nil
	} else if result.Err() != nil {
		t.Fatal(result.Err())
		return
	}
	_ = result
}

func TestMYSQL(t *testing.T) {
	var (
		table = "test"
	)

	db, err := newMYSQL()
	if err != nil {
		t.Fatal(err)
		return
	}
	builder := &MYSQL{}

	var expr Expressions = AND{}
	expr = expr.Set(
		"",
		IN{Column: "id", Values: []interface{}{"1", "2"}},
		LIKE{Column: "name", Values: "alice"},
	)

	expr.Build(builder)
	result := make(map[string]interface{})
	err = db.Table(table).
		Where(builder.SQL.String(), builder.Vars...).
		Find(&result).
		Error
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestClause(t *testing.T) {
	clause := New()
	exprIN1, err := clause.GetExpression("in", "id", "1", "2")
	if err != nil {
		t.Fatal(err)
		return
	}

	exprIN2, err := clause.GetExpression("in", "name", "alice")
	if err != nil {
		t.Fatal(err)
		return
	}
	expr, err := clause.GetExpression(
		"and",
		"",
		exprIN1, exprIN2,
	)
	if err != nil {
		t.Fatal(err)
		return
	}

	_ = expr
}
