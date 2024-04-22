package tasquelite

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormTaskStorage[D interface{}] struct {
	db *gorm.DB
}

func NewGormTaskStorageManager[D interface{}](dbName string, model *D) (*GormTaskStorage[D], error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema for the EmailTask model.
	err = db.AutoMigrate(model)
	if err != nil {
		return nil, err
	}

	return &GormTaskStorage[D]{db: db}, nil
}

// SaveTaskToStorage saves a task to the storage.
func (gts *GormTaskStorage[D]) SaveTaskToStorage(ctx context.Context, task *D) error {
	return gts.db.WithContext(ctx).Create(task).Error
}

// GetTasksFromStorage retrieves tasks from the storage.
func (gts *GormTaskStorage[D]) GetTasksFromStorage(ctx context.Context, limit int) ([]D, error) {
	var tasks []D
	model := new(D)
	err := gts.db.WithContext(ctx).Model(model).Limit(limit).Find(&tasks).Error
	return tasks, err
}

// DeleteTaskFromStorage deletes a task from the storage.
func (gts *GormTaskStorage[D]) DeleteTaskFromStorage(ctx context.Context, task *D) error {
	model := new(D)
	return gts.db.WithContext(ctx).Model(model).Delete(task).Where(task).Error
}
