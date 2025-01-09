package dao

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	defaultPatchSize = 100
)

var ErrRecordNotFound = errors.New("未查到相关信息！")
var ErrExistSameName = errors.New("name exist")

type Store struct {
	db *gorm.DB
}

func New(ctx *gin.Context, db *gorm.DB) *Store {
	s := &Store{
		db: db,
	}
	return s
}

func (s *Store) delete(model any) error {
	return s.db.Delete(model).Error
}

func (s *Store) create(model any) error {
	return s.db.Create(model).Error
}

func (s *Store) save(model any) error {
	return s.db.Save(model).Error
}

func (s *Store) update(model any, attrs ...any) error {
	return s.db.Model(model).UpdateColumns(attrs).Error
}

func (s *Store) updates(model any, values any) error {
	return s.db.Model(model).Updates(values).Error
}

func (s *Store) Tx(fn func(txDB *gorm.DB) error) (err error) {
	txDB := s.db.Begin()
	if err = fn(txDB); err != nil {
		txDB.Rollback()
		return err
	}

	if err = txDB.Commit().Error; err != nil {
		txDB.Rollback()
		return err
	}
	return nil
}

func (s *Store) createInBatches(value any, batchSize int) error {
	return s.db.CreateInBatches(value, batchSize).Error
}

func (s *Store) like(query map[string]string) *gorm.DB {
	// 动态构建 WHERE 条件
	dbQuery := s.db
	for key, value := range query {
		if value == "" {
			continue
		}
		likeQuery := "%" + value + "%"
		dbQuery = dbQuery.Where(fmt.Sprintf("%s LIKE ?", key), likeQuery)
	}
	return dbQuery
}
