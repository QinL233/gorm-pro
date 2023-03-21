package dao

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func OneKey[T any](db *gorm.DB, key any) (T, error) {
	var result T
	//根据id查询时使用first默认以id排序
	if err := db.Model(&result).First(&result, key).Error; err != nil {
		return result, err
	}
	return result, nil
}

func OneKeyTry[T any](db *gorm.DB, key any) (T, error) {
	if result, err := OneKey[T](db, key); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return result, err
	} else {
		//注意泛型返回时必定会返回一个属性为zero的结构，即使用时执行通过属性判断
		return result, nil
	}
}

func One[T any](db *gorm.DB, condition interface{}, args ...interface{}) (T, error) {
	var result T
	if err := db.Model(&result).Where(condition, args...).First(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

func OneTry[T any](db *gorm.DB, condition interface{}, args ...interface{}) (T, error) {
	if result, err := One[T](db, condition, args...); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return result, err
	} else {
		//注意泛型返回时必定会返回一个属性为zero的结构，即使用时执行通过属性判断
		return result, nil
	}
}

func OneEntity[T any](db *gorm.DB, entity T) (T, error) {
	var result T
	if err := db.Model(&entity).Where(&entity).Take(&result).Error; err != nil {
		return result, err
	} else {
		return result, nil
	}
}

func OneEntityTry[T any](db *gorm.DB, entity T) (T, error) {
	if result, err := OneEntity[T](db, entity); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return result, err
	} else {
		//注意泛型返回时必定会返回一个属性为zero的结构，即使用时执行通过属性判断
		return result, nil
	}
}

func OneScope[T any](db *gorm.DB, scope func(db *gorm.DB) *gorm.DB) (T, error) {
	var result T
	if err := db.Model(&result).Scopes(scope).First(&result).Error; err != nil {
		return result, err
	}
	return result, nil
}

func OneScopeTry[T any](db *gorm.DB, scope func(db *gorm.DB) *gorm.DB) (T, error) {
	if result, err := OneScope[T](db, scope); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return result, err
	} else {
		//注意泛型返回时必定会返回一个属性为zero的结构，即使用时执行通过属性判断
		return result, nil
	}
}
