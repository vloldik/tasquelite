package tasquelite

import (
	"context"

	"github.com/vloldik/tasque"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormTaskStorage[T tasque.Task] struct {
	db        *gorm.DB
	batchSize int
}

func NewGormTaskStorageManager[T tasque.Task](dbName string, model *T, batchSize int) (*GormTaskStorage[T], error) {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema for the EmailTask model.
	err = db.AutoMigrate(model)
	if err != nil {
		return nil, err
	}

	return &GormTaskStorage[T]{db: db, batchSize: batchSize}, nil
}

// SaveTaskToStorage saves a task to the storage.
func (gts *GormTaskStorage[T]) SaveTaskToStorage(ctx context.Context, task *T) error {
	return gts.db.WithContext(ctx).Create(task).Error
}

// GetTasksFromStorage retrieves tasks from the storage.
func (gts *GormTaskStorage[T]) GetTasksFromStorage(ctx context.Context) ([]T, error) {
	var tasks []T
	model := new(T)
	err := gts.db.WithContext(ctx).Model(model).Limit(gts.TaskFromStorageBatchCount()).Find(&tasks).Error
	return tasks, err
}

// DeleteTaskFromStorage deletes a task from the storage.
func (gts *GormTaskStorage[T]) DeleteTaskFromStorage(ctx context.Context, task *T) error {
	model := new(T)
	return gts.db.WithContext(ctx).Model(model).Delete(task).Where(task).Error
}

// TaskFromStorageBatchCount returns the batch count for getting tasks from storage.
func (gts *GormTaskStorage[T]) TaskFromStorageBatchCount() int {
	// Implement the logic to determine the batch count.
	return gts.batchSize
}
