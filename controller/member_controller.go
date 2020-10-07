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
	"log"
	"strconv"
)

// 根据项目id发送成员数组
func SendProjectMembers(ctx *gin.Context) {
	DB = common.GetDB()
	// 获取项目名字
	projectName := ctx.Param("name")
	project, err := service.GetProjectByName(projectName)
	if !util.IsOk(err) {
		fmt.Println(err.Message)
		response.Fail(ctx, nil, err.Message)
		return
	}
	projectId := project.Id

	// 获取项目的所有成员
	users, err1 := service.GetAllMembersByProjectId(projectId)
	if !util.IsOk(err1) {
		response.Fail(ctx, nil, "获取成员数组失败！用户不存在！")
		return
	}
	fmt.Printf("%+v", *users)
	count := len(*users)
	response.Success(ctx, gin.H{
		"Count":   count,
		"Members": *users,
	}, "获取成员列表成功！")
}

// 根据项目id和成员email，或者成员的mixin ID为项目添加成员
func AddProjectMember(ctx *gin.Context) {
	DB := common.GetDB()
	// 解析获取项目id和成员email, ctx不能被解析两次
	var infoMap = make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&infoMap)
	projectId, _ := service.GetProjectIdByName(ctx.Param("name"))
	memberEmail := infoMap["memberEmail"]
	mixinId := infoMap["mixinId"]

	if len(memberEmail) == 0 && len(mixinId) == 0 {
		response.Fail(ctx, nil, "成员email或者mixin ID获取错误！")
		log.Panicln("成员email或者mixin ID获取错误！")
		return
	}

	var user *model.User
	var err1 *util.Err
	if len(memberEmail) != 0 { // 由成员email获取用户
		user, err1 = service.GetUserByEmail(memberEmail)
		if !util.IsOk(err1) {
			response.Fail(ctx, nil, err1.Message)
			return
		}
	} else { // 由mixin ID 获取用户
		user, err1 = service.GetUserByMixinId(mixinId)
		if !util.IsOk(err1) {
			response.Fail(ctx, nil, err1.Message)
			return
		}
	}

	// 在members表中添加成员
	var member model.Member
	userId := user.Id
	member.ProjectId = projectId
	member.UserId = userId
	if err := DB.Create(&member).Error; err != nil {
		response.Fail(ctx, nil, user.Email+"该成员已经存在！")
		return
	}
	// 在member_wallets中添加
	if err := service.CreateMemberWallet(projectId, userId); !util.IsOk(err) {
		response.Fail(ctx, nil, "用户钱包创建错误！")
		return
	}

	response.Success(ctx, nil, "添加成员成功！")
}

// 为项目添加多个成员
func AddProjectMembers(ctx *gin.Context) {
	DB := common.GetDB()
	projectName := ctx.Param("name")
	projectId, _ := service.GetProjectIdByName(projectName)
	var infoMap = make(map[string]interface{})
	json.NewDecoder(ctx.Request.Body).Decode(&infoMap)
	members := infoMap["members"]

	// 添加成员
	if err := service.CreateMembers(members, projectName); !util.IsOk(err) {
		response.Fail(ctx, nil, "添加成员失败！")
		log.Panicln(err)
	}

	//distribution := service.GetDistributionByProjectId(projectId)
	// 添加成员钱包
	if err := service.CreateMemberWallets(DB, projectId); !util.IsOk(err) {
		response.Fail(ctx, nil, "添加成员钱包时出错！")
		log.Panicln(err)
	}
	response.Success(ctx, nil, "添加成员成功！")
}

// 根据成员id为项目删除成员
func DeleteProjectMember(ctx *gin.Context) {

	// 解析获取成员id
	var infoMap = make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&infoMap)
	userId, err := strconv.Atoi(infoMap["userId"])
	if err != nil || userId == 0 {
		fmt.Println("member_controller/DeleteProjectMember: "+
			"成员id获取错误！", err)
		response.Fail(ctx, nil, "成员id获取错误！")
		return
	}

	projectId, _ := service.GetProjectIdByName(ctx.Param("name"))

	// 删除成员
	err1 := service.DeleteMember(int64(userId), projectId)
	if !util.IsOk(err1) {
		response.Fail(ctx, nil, err1.Message)
		return
	}
	response.Success(ctx, nil, "删除成功！")

}
