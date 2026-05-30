package repository

import "gorm.io/gorm"

// ApplyIncludeDeleted 参数控制是否包含软删除记录
func ApplyIncludeDeleted(q *gorm.DB, includeDeleted bool) *gorm.DB {
	if includeDeleted { return q.Unscoped() }
	return q
}

// ForceDelete 硬删除（绕过 GORM 软删除）
func ForceDelete(db *gorm.DB, model interface{}, id string) error {
	return db.Unscoped().Where("id = ?", id).Delete(model).Error
}

// Restore 恢复软删除记录
func Restore(db *gorm.DB, model interface{}, id string) error {
	return db.Unscoped().Model(model).Where("id = ?", id).Update("deleted_at", nil).Error
}
