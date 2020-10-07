package controller

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/response"
	"claps-admin/service"
	"claps-admin/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
)

var DB *gorm.DB

// 从数据库读取管理员信息并返回给前端
func SendAdminTable(ctx *gin.Context) {
	DB = common.GetDB()
	admins, count := service.GetAllAdmins()
	response.Success(ctx, gin.H{
		"Count":  count,
		"Admins": admins,
	}, "返回用户信息成功！")
}

// 向数据库添加管理员
func AddNewAdmin(ctx *gin.Context) {
	DB = common.GetDB()
	var adminMap = make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&adminMap)
	// 指针类型
	var newAdmin *model.Admin
	var err *util.Err
	if newAdmin, err = service.CreateAdmin(adminMap); err.Code != util.OK {
		response.Fail(ctx, nil, err.Message)
		return
	}
	fmt.Println("创建的newAdmin为：", *newAdmin)
	DB.Create(newAdmin)

	response.Success(ctx, nil, "添加管理员成功！")
}

// 向数据库中添加多个管理员
func AddNewAdmins(ctx *gin.Context) {
	var infoMap = make(map[string]interface{})
	json.NewDecoder(ctx.Request.Body).Decode(&infoMap)
	admins := infoMap["admins"]

	// 添加管理员
	if err := service.CreateAdmins(admins); !util.IsOk(err) {
		response.Fail(ctx, nil, "添加管理员失败！")
		log.Panicln(err)
	}

	response.Success(ctx, nil, "添加管理员成功！")
}

// 从数据库删除管理员
func DeleteAdmin(ctx *gin.Context) {
	DB = common.GetDB()
	// 获取账号
	var infoMap = make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&infoMap)
	email := infoMap["email"]

	// 删除管理员
	if err := service.DeleteAdminByAccount(email); !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		log.Panicln("删除管理员失败！")
		return
	}
	response.Success(ctx, nil, "删除成功！")
}
