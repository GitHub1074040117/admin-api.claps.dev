package controller

import (
	"claps-admin/response"
	"claps-admin/service"
	"claps-admin/util"
	"github.com/gin-gonic/gin"
	"log"
)

// 根据获取到的用户id发送用户信息
func SendUserInfo(ctx *gin.Context) {
	// 获取用户id
	userId, _ := service.GetUserIdByName(ctx.Param("name"))

	// 根据用户id获取用户
	user, err1 := service.GetUserById(userId)
	if !util.IsOk(err1) {
		response.Fail(ctx, nil, "用户信息获取失败！")
		log.Panicln("用户信息获取失败！")
	}
	response.Success(ctx, gin.H{
		"User": *user,
	}, "用户信息获取成功")
}

func SendAllUsers(ctx *gin.Context) {
	users, err := service.GetAllUsers()
	if !util.IsOk(err) {
		response.Fail(ctx, nil, "获取用户列表失败！")
		return
	}
	response.Success(ctx, gin.H{
		"Users": users,
	}, "返回用户列表成功！")

}
