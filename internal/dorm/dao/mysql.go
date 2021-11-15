package dao

import (
	"context"
	"strings"

	"git.internal.yunify.com/qxp/molecule/internal/dorm/clause"
	"gorm.io/gorm"
)

// MYSQL mysql dao
type MYSQL struct {
	DB *gorm.DB
}

// FindOne find one entity
func (m *MYSQL) FindOne(ctx context.Context, builder clause.Builder) (Data, error) {
	mysqlBuilder, err := instantiationMYSQL(builder)
	if err != nil {
		return nil, err
	}
	var result = make(Data)
	err = m.DB.
		Where(mysqlBuilder.SQL, mysqlBuilder.Vars...).
		Find(&result).
		Error
	return result, err
}

// Find find entities
func (m *MYSQL) Find(ctx context.Context, builder clause.Builder, findOpt FindOptions) ([]Data, error) {
	mysqlBuilder, err := instantiationMYSQL(builder)
	if err != nil {
		return nil, err
	}

	var result = make([]Data, 0)
	err = m.DB.
		Where(mysqlBuilder.SQL, mysqlBuilder.Vars...).
		Offset(int((findOpt.Page - 1) * findOpt.Size)).
		Limit(int(findOpt.Size)).
		Order(mysqlSort(findOpt.Sort...)).
		Find(&result).
		Error

	return result, err
}

// Count count entities
func (m *MYSQL) Count(ctx context.Context, builder clause.Builder) (int64, error) {
	mysqlBuilder, err := instantiationMYSQL(builder)
	if err != nil {
		return 0, err
	}
	var total int64
	err = m.DB.
		Where(mysqlBuilder.SQL, mysqlBuilder.Vars...).
		Count(&total).
		Error
	return total, err
}

// Insert insert entities
func (m *MYSQL) Insert(ctx context.Context, entity ...interface{}) error {
	return m.DB.
		CreateInBatches(entity, len(entity)).
		Error
}

// Update update entities
func (m *MYSQL) Update(ctx context.Context, builder clause.Builder, entity interface{}) (int64, error) {
	mysqlBuilder, err := instantiationMYSQL(builder)
	if err != nil {
		return 0, err
	}
	tx := m.DB.Where(mysqlBuilder.SQL, mysqlBuilder.Vars).Updates(entity)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}

// Delete delete entities with condition
func (m *MYSQL) Delete(ctx context.Context, builder clause.Builder) (int64, error) {
	mysqlBuilder, err := instantiationMYSQL(builder)
	if err != nil {
		return 0, err
	}
	tx := m.DB.Where(mysqlBuilder.SQL, mysqlBuilder.Vars...).Delete(map[string]interface{}{})
	if tx.Error != nil {
		return 0, tx.Error
	}
	return tx.RowsAffected, nil
}
func instantiationMYSQL(builder clause.Builder) (*clause.MYSQL, error) {
	mysql, ok := builder.(*clause.MYSQL)
	if !ok {
		return nil, ErrAssertBuilder
	}
	return mysql, nil
}

func mysqlSort(array ...string) string {
	sort := make([]string, 0, len(array))
	for _, elem := range array {
		if strings.HasPrefix(elem, "-") {
			sort = append(sort, elem[1:]+" DESC")
			continue
		}
		sort = append(sort, elem+" ASC")
	}
	return strings.Join(sort, ",")
}
