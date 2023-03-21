package dao

import "gorm.io/gorm"

/**
指定参数进行查询
*/
func List[T any](db *gorm.DB, condition interface{}, args ...interface{}) ([]T, error) {
	return ListSortTo[T, T](db, "", condition, args...)
}

func ListSort[T any](db *gorm.DB, sort string, condition interface{}, args ...interface{}) ([]T, error) {
	return ListSortTo[T, T](db, sort, condition, args...)
}

func ListSortTo[T any, E any](db *gorm.DB, sort string, condition interface{}, args ...interface{}) ([]E, error) {
	return ListFieldSortTo[T, E](db, nil, sort, condition, args)
}

func ListFieldSortTo[T any, E any](db *gorm.DB, field []string, sort string, condition interface{}, args ...interface{}) ([]E, error) {
	var result []E
	var entity T
	query := db.Model(&entity).Where(condition, args...)
	if len(field) > 0 {
		query.Select(field)
	}
	if sort != "" {
		query.Order(sort)
	}
	if err := query.Find(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

/**
指定对象进行查询
*/
func ListEntity[T any](db *gorm.DB, entity T) ([]T, error) {
	return ListEntitySortTo[T, T](db, "", entity)
}

func ListEntitySort[T any](db *gorm.DB, sort string, entity T) ([]T, error) {
	return ListEntitySortTo[T, T](db, sort, entity)
}

func ListEntitySortTo[T any, E any](db *gorm.DB, sort string, entity T) ([]E, error) {
	return ListEntityFieldSortTo[T, E](db, nil, sort, entity)
}

func ListEntityFieldSortTo[T any, E any](db *gorm.DB, field []string, sort string, entity T) ([]E, error) {
	var result []E
	query := db.Model(&entity).Where(&entity)
	if len(field) > 0 {
		query.Select(field)
	}
	if sort != "" {
		query.Order(sort)
	}
	if err := query.Find(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

/**
自定义scope进行查询
*/
func ListScope[T any](db *gorm.DB, scope func(db *gorm.DB) *gorm.DB) ([]T, error) {
	return ListScopeSortTo[T, T](db, "", scope)
}

func ListScopeSort[T any](db *gorm.DB, sort string, scope func(db *gorm.DB) *gorm.DB) ([]T, error) {
	return ListScopeSortTo[T, T](db, sort, scope)
}

func ListScopeSortTo[T any, E any](db *gorm.DB, sort string, scope func(db *gorm.DB) *gorm.DB) ([]E, error) {
	return ListScopeFieldSortTo[T, E](db, nil, sort, scope)
}

func ListScopeFieldSortTo[T any, E any](db *gorm.DB, field []string, sort string, scope func(db *gorm.DB) *gorm.DB) ([]E, error) {
	var result []E
	var entity T
	query := db.Model(&entity).Scopes(scope)
	if len(field) > 0 {
		query.Select(field)
	}
	if sort != "" {
		query.Order(sort)
	}
	if err := query.Find(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}
