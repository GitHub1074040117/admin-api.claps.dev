package common

import (
	"claps-admin/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	//配置项读取
	driverName := viper.GetString("datasource.driverName")

	args := fmt.Sprintf("host=%s port=%s dbname=%s user=%s  password=%s sslmode=disable",
		viper.GetString("datasource.host"),
		viper.GetString("datasource.port"),
		viper.GetString("datasource.database"),
		viper.GetString("datasource.username"),
		viper.GetString("datasource.password"))

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}
	db.SingularTable(true)
	DB = db
	// 数据库初始化表格
	if !db.HasTable(&model.User{}) {
		db.AutoMigrate(&model.User{})
	}
	if !db.HasTable(&model.Admin{}) {
		db.AutoMigrate(&model.Admin{})
	}
	if !db.HasTable(&model.Project{}) {
		db.AutoMigrate(&model.Project{})
	}
	if !db.HasTable(&model.ProjectToMericoGroup{}) {
		db.AutoMigrate(&model.ProjectToMericoGroup{})
	}
	if !db.HasTable(&model.Member{}) {
		db.AutoMigrate(&model.Member{})
	}
	if !db.HasTable(&model.MemberWallet{}) {
		db.AutoMigrate(&model.MemberWallet{})
	}
	if !db.HasTable(&model.Repository{}) {
		db.AutoMigrate(&model.Repository{})
	}
	if !db.HasTable(&model.RepositoryToMericoProject{}) {
		db.AutoMigrate(&model.RepositoryToMericoProject{})
	}
	if !db.HasTable(&model.Wallet{}) {
		db.AutoMigrate(&model.Wallet{})
	}
	if !db.HasTable(&model.Bot{}) {
		db.AutoMigrate(&model.Bot{})
	}
	if !db.HasTable(&model.Transaction{}) {
		db.AutoMigrate(&model.Transaction{})
	}
	if !db.HasTable(&model.Transfer{}) {
		db.AutoMigrate(&model.Transfer{})
	}
	if !db.HasTable(&model.Asset{}) {
		db.AutoMigrate(&model.Asset{})
	}

	return db
}

func GetDB() *gorm.DB {
	return DB
}
