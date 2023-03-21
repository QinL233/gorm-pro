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

func TestList(t *testing.T) {
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
	result1, _ := dao.ListEntityFieldSortTo[User, struct{ Username, Password string }](db, []string{"username", "password"}, "id", User{Id: 1})
	fmt.Println(result1)
	result2, _ := dao.ListEntitySortTo[User, User](db, "id", User{Id: 1})
	fmt.Println(result2)
}
