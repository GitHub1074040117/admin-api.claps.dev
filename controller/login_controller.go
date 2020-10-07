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
	"github.com/spf13/viper"
	"log"
)

type Conf struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

var conf = Conf{}
var admin = model.Admin{}

const AuthOff = "off"

// 返回给前端登录github信息
func SendAuthUrl(ctx *gin.Context) {

	// 从配置文件中读取
	conf.ClientId = viper.GetString("oauthConf.clientId")
	conf.ClientSecret = viper.GetString("oauthConf.clientSecret")
	conf.RedirectUrl = viper.GetString("oauthConf.redirectUrl")

	// 返回结果到前端
	response.Success(ctx, gin.H{
		// github授权界面url
		"AuthorizeUrl": "https://github.com/login/oauth/authorize?client_id=" + conf.ClientId + "&redirect_uri=" + conf.RedirectUrl,
	}, "获取client数据成功！")
}

// 获取Token，并检查用户是否为管理员，是则允许登录主界面
func CheckLoginCode(ctx *gin.Context) {
	var requestMap = make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&requestMap)
	var code = requestMap["Code"]
	if len(code) == 0 {
		response.Fail(ctx, nil, "获取code失败")
		log.Panicln("出错了，获取的code长度为0")
		return
	}
	code = code[6:]

	authUrl := service.GetTokenAuthUrl(conf.ClientId, conf.ClientSecret, code)
	token, err := service.GetToken(authUrl)
	if !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		log.Panicln(err.Message)
	}

	user, err := service.GetUserInfo(token)
	if !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		log.Panicln(err)
	}

	admin.AvatarUrl = *user.AvatarURL
	admin.Name = *user.Name
	admin.Email = *user.Login
	// 判断是否为高级管理员
	if admin.Email == viper.GetString("superAdmin.account") {
		admin.Role = "super"
	} else {
		admin.Role = "common"
	}

	DB := common.GetDB()
	// 不需要验证
	if viper.GetString("authOff") == AuthOff {
		response.Success(ctx, nil, "用户身份验证成功！")
		return
	}
	// 查询数据库,判断当前登录的用户是否为管理员
	if !service.IsAdminExist(admin.Email) {
		response.Fail(ctx, nil, "该用户不是管理员")
		fmt.Println("该用户不是管理员！")
		return
	} else { // 更新用户在数据库中的信息
		service.UpdateUserInfo(DB, admin)
	}
	response.Success(ctx, nil, "用户身份验证成功！")
}

// 返回用户信息和本平台的Token给前端
func SendLoginInfo(ctx *gin.Context) {
	// 发放token
	token, err := common.ReleaseToken(admin)
	if !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		return
	}
	response.Success(ctx, gin.H{
		"Token": token,
		"User":  admin,
	}, "返回用户信息成功！")
}
