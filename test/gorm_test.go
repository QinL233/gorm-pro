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
	driver().Create([]User{
		{1, "admin1", "admin1"},
		{2, "admin2", "admin2"},
	})
	driver().Save([]User{
		{1, "root1", "root1"},
		{2, "root2", "root2"},
	})
	driver().Save([]User{
		{Username: "admin1", Password: "admin1"},
		{Username: "admin2", Password: "admin2"},
	})
	//save变更会导致所有列变更
	driver().Save([]User{
		{3, "", "admin"},
		{4, "", "admin"},
	})
}

func TestDel(t *testing.T) {
	dao.RemoveKey[User](driver(), 1)
	dao.Remove[User](driver(), "id = ?", 2)
	//entity传递泛型可以在参数中
	dao.RemoveEntity(driver(), User{Id: 3})
	dao.RemoveScope[User](driver(), func(db *gorm.DB) *gorm.DB {
		return db.Where("id in ?", []int{1, 2, 3, 4})
	})
}

func TestUpdate(t *testing.T) {
	dao.Update[User](driver(), "username", "admin3", "id = ?", 3)
	dao.UpdateEntity(driver(), "username", gorm.Expr("id + ?", 1), User{Id: 4})
	dao.UpdateScope[User](driver(), "username", "test", func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", 4)
	})
	dao.Updates[User](driver(), map[string]interface{}{
		"username": "admin33",
		"password": "admin33",
	}, "id = ?", 3)
	dao.UpdatesEntity(driver(), map[string]interface{}{
		"username": gorm.Expr("id + ?", 1),
		"password": gorm.Expr("id + ?", 1),
	}, User{Id: 4})
	dao.UpdatesScope[User](driver(), map[string]interface{}{
		"username": "test",
		"password": "test",
	}, func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", 4)
	})
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
		"id desc", 10,
		"username = ?", "root")
	fmt.Println(list1)
	list2, _ := dao.ListEntitySortLimitFieldTo[User, User](
		driver(),
		nil,
		"id asc", 10,
		User{Id: 1})
	fmt.Println(list2)
	list3, _ := dao.ListScopeTo[User, User](driver(), func(db *gorm.DB) *gorm.DB {
		return db.Select("username").Limit(10)
	})
	fmt.Println(list3)
}

func TestPage(t *testing.T) {
	count1, page1, _ := dao.PageSortFieldTo[User, struct{ Username string }](
		driver(),
		[]string{"username"},
		"id asc", 10, 1,
		"id > ?", 0)
	fmt.Println(count1, page1)
	count2, page2, _ := dao.PageEntitySortFieldTo[User, struct{ Username string }](
		driver(),
		[]string{"username"},
		"id asc", 10, 1,
		User{Username: "root"})
	fmt.Println(count2, page2)
	count3, page3, _ := dao.PageScopeTo[User, struct{ Username string }](
		driver(), 10, 1,
		func(db *gorm.DB) *gorm.DB {
			//page函数使用了select count(*)因此注意scope不可以覆盖select
			//db.Select("username")
			return db.Where("id > ?", 1)
		},
	)
	fmt.Println(count3, page3)
}
