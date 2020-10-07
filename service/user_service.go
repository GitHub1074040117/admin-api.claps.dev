package service

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/util"
	"log"
)

// 按邮箱查找数据库中的user表，并返回user指针
func GetUserByEmail(email string) (*model.User, *util.Err) {
	var user model.User
	db := common.GetDB()
	db.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		log.Panicln("获取用户失败！用户email不存在！")
	}
	return &user, util.Success()
}

// 通过user_id查找数据库的user表，返回指针
func GetUserById(userId int64) (*model.User, *util.Err) {
	var user model.User
	db := common.GetDB()
	db.Where("id = ?", userId).First(&user)
	if user.Id == 0 {
		log.Panicln("获取用户失败！用户id不存在！id = ", userId)
	}
	return &user, util.Success()
}

// 按mixin ID 查数据库的user表，并返回user指针
func GetUserByMixinId(mixinId string) (*model.User, *util.Err) {
	var user model.User
	db := common.GetDB()
	db.Where("mixin_id = ?", mixinId).First(&user)
	if user.Id == 0 {
		log.Println("获取用户失败！用户mixin ID不存在！", mixinId)
		return nil, util.Fail("获取用户失败！用户mixin ID不存在！")
	}
	return &user, util.Success()
}

func GetUserIdByName(userName string) (int64, *util.Err) {
	DB := common.GetDB()
	var user model.User
	if err := DB.Where(" name = ? ", userName).First(&user).Error; err != nil {
		log.Panicln(err)
	}
	if user.Id == 0 {
		log.Panicln("查找不到该用户id", userName)
	}
	return user.Id, util.Success()
}

func GetAllUsers() (*[]model.User, *util.Err) {
	DB := common.GetDB()
	var users []model.User
	err := DB.Find(&users).Error
	if err != nil {
		log.Panicln("获取用户失败！", err)
	}
	return &users, util.Success()
}

/*func GetNotProjectUsers(projectName string) (*[]model.User, *util.Err) {
	DB := common.GetDB()
	var users []model.User
	err := DB.Where("  ")
}*/
