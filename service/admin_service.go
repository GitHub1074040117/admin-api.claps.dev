package service

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/util"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func GetAllAdmins() (*[]model.Admin, int) {
	DB := common.GetDB()
	var admins []model.Admin
	var count int
	// 读取长度
	DB.Table("admin").Count(&count)
	// 读取所有记录
	DB.Find(&admins)
	return &admins, count
}

// 在数据库中判断一名申请登录的用户是否为管理员
func IsAdminExist(email string) bool {
	db := common.GetDB()
	var admin model.Admin
	fmt.Println("待查找的管理员账号：", email)
	db.Where("email = ?", email).First(&admin)
	fmt.Println("找到的管理员ID：", admin.UserId)
	if admin.UserId != 0 {
		return true
	}
	return false
}

// 通过map创建管理员, 返回指针类型
func CreateAdmin(adminMap map[string]string) (*model.Admin, *util.Err) {
	newAdmin := model.Admin{}
	fmt.Println("（service/CreateAdmin）map:", adminMap)
	// 判断账号是否为空
	if len(adminMap["account"]) == 0 {
		return &newAdmin, util.Fail("出错了，获取的账号为空！")
	}
	newAdmin.Email = adminMap["account"]
	//判断数据库中是否存在
	if IsAdminExist(newAdmin.Email) {
		return &newAdmin, util.Fail("创建失败！该用户已存在！")
	}

	// 判断哪个值为空则赋值
	if len(adminMap["name"]) == 0 {
		newAdmin.Name = util.RandomString(5)
	} else {
		newAdmin.Name = adminMap["name"]
	}
	if len(adminMap["role"]) == 0 {
		newAdmin.Role = "common"
	} else {
		newAdmin.Role = adminMap["common"]
	}
	if len(adminMap["avatarUrl"]) == 0 {
		newAdmin.AvatarUrl = viper.GetString("superAdmin.avatarUrl")
	} else {
		newAdmin.AvatarUrl = adminMap["avatarUrl"]
	}
	return &newAdmin, util.Success()
}

// 同时添加多个管理员
func CreateAdmins(admins interface{}) *util.Err {
	DB := common.GetDB()
	var err *util.Err
	for _, userName := range admins.([]interface{}) {
		var admin model.Admin
		var user *model.User
		admin.UserId, err = GetUserIdByName(userName.(string))
		if !util.IsOk(err) {
			log.Panicln("向数据库添加管理员失败！")
		}
		if user, err = GetUserById(admin.UserId); !util.IsOk(err) {
			log.Panicln("根据用户id获取用户失败！")
		}
		admin.Name = user.Name
		admin.MixinId = user.MixinId
		admin.AvatarUrl = user.AvatarUrl
		admin.Role = model.Common
		admin.Email = user.Email
		if err1 := DB.Create(&admin).Error; err1 != nil {
			log.Panicln(err1)
		}
	}
	return util.Success()
}

// 根据账号删除管理员
func DeleteAdminByAccount(email string) *util.Err {
	db := common.GetDB()
	// 检查管理员是否存在
	if !IsAdminExist(email) {
		log.Panicln("删除失败！管理员不在数据库！")
	}
	// 删除
	db.Where("email = ?", email).Delete(&model.Admin{})
	return util.Success()
}

// 在数据库中初始化一名高级管理员
func InitAdmin() {
	DB := common.GetDB()
	superAdmin := model.Admin{
		UserId:    viper.GetInt64("superAdmin.userId"),
		Name:      viper.GetString("superAdmin.name"),
		Email:     viper.GetString("superAdmin.account"),
		AvatarUrl: viper.GetString("superAdmin.avatarUrl"),
		Role:      viper.GetString("superAdmin.role"),
	}

	if IsAdminExist(superAdmin.Email) {
		fmt.Println("高级用户已经存在！")
		DB.Model(&superAdmin).Update("role", model.Super)
	} else {
		fmt.Println("高级用户不存在！已经创建")
		DB.Create(&superAdmin)
	}
}
