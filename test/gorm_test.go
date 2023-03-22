package test

import (
	"fmt"
	"github.com/QinL233/gorm-pro/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"testing"
)

type User struct {
	Id       int    `gorm:"primaryKey" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func driver() *gorm.DB {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",
		"root",
		"172.17.181.4",
		3306,
		"test")),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), //日志级别
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, //取消表明被加s
			},
			DisableForeignKeyConstraintWhenMigrating: true, //取消外键约束
			SkipDefaultTransaction:                   true, //禁用默认事务可以提升性能
		})
	if err != nil {
		panic(err)
	}
	return db
}

func TestSaveBatch(t *testing.T) {

}

func TestDel(t *testing.T) {
}

func TestUpdate(t *testing.T) {
}

func TestCount(t *testing.T) {
	count1, _ := dao.Count[User](driver(), "username in ?", []string{"root", "admin"})
	fmt.Println(count1)
	count2, _ := dao.CountEntity[User](driver(), User{Password: "root"})
	fmt.Println(count2)
	count3, _ := dao.CountScope[User](driver(), func(db *gorm.DB) *gorm.DB {
		return db.Where("username like ?", "%oo%")
	})
	fmt.Println(count3)

}

func TestOne(t *testing.T) {
	one1, _ := dao.OneKeyFieldTo[User, struct{ Token string }](driver(), []string{"username as token"}, 2)
	fmt.Println(one1)
	one2, _ := dao.OneFieldTo[User, struct{ Token string }](driver(), []string{"username as token"}, "password = ?", "admin")
	fmt.Println(one2)
	one3, _ := dao.OneEntityFieldTo[User, struct{ Token string }](driver(), []string{"username as token"}, User{Username: "root"})
	fmt.Println(one3)
	one4, _ := dao.OneScopeTo[User, struct{ Token string }](driver(), func(db *gorm.DB) *gorm.DB {
		return db.Select("concat(username,password) as token")
	})
	fmt.Println(one4)
}

func TestList(t *testing.T) {
	list1, _ := dao.ListSortLimitFieldTo[User, struct{ Username, Password string }](
		driver(),
		[]string{"username", "password"},
		"id", 10,
		"username = ?", "root")
	fmt.Println(list1)
	list2, _ := dao.ListEntitySortLimitFieldTo[User, User](
		driver(),
		nil,
		"id", 10,
		User{Id: 1})
	fmt.Println(list2)
	list3, _ := dao.ListScopeTo[User, User](driver(), func(db *gorm.DB) *gorm.DB {
		return db.Select("username").Limit(10)
	})
	fmt.Println(list3)
}
