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
)

/*
util.Err.Message是需要在controller打印出来的错误信息,
*/

// 添加新项目
func AddNewProject(ctx *gin.Context) {
	DB := common.GetDB()
	var err *util.Err
	var projectId int64

	var infoMap = make(map[string]interface{})

	json.NewDecoder(ctx.Request.Body).Decode(&infoMap)

	//distribution := infoMap["distribution"].(string)
	repoType := infoMap["repoType"].(string)
	repoUrl := infoMap["repoUrl"].(string)
	repoName := infoMap["repoName"].(string)
	projectName := infoMap["projectName"].(string)
	description := infoMap["description"].(string)
	avatarUrl := infoMap["avatarUrl"].(string)
	members := infoMap["members"]

	// 添加项目
	if projectId, err = service.CreateProject(projectName, description, avatarUrl); !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		log.Panicln(err)
	}

	// 添加成员
	if err = service.CreateMembers(members, projectName); !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		// 添加失败删除项目
		service.DeleteProjectById(projectId)
		log.Panicln(err.Message)
	}

	// 添加仓库
	if _, err = service.CreateRepository(projectId, repoType, repoUrl, repoName); !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		// 添加失败删除项目
		service.DeleteProjectById(projectId)
		log.Panicln(err.Message)
	}

	// 添加机器人
	if err = service.Create4BotsForProject(projectId); !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		// 添加失败删除项目
		service.DeleteProjectById(projectId)
		return
	}

	// 添加项目钱包
	if err = service.CreateProjectWallets(DB, projectId); !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		// 添加失败删除项目
		service.DeleteProjectById(projectId)
		return
	}

	// 添加成员钱包
	if err = service.CreateMemberWallets(DB, projectId); !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		// 添加失败删除项目
		service.DeleteProjectById(projectId)
		return
	}

	response.Success(ctx, nil, "项目添加成功！")
}

// 更新项目信息
func UpdateProjectInfo(ctx *gin.Context) {
	infoMap := make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&infoMap)
	projectName := ctx.Param("name")
	editedName := infoMap["editedName"]
	avatarUrl := infoMap["avatarUrl"]
	desc := infoMap["description"]
	fmt.Println("avatar:" + avatarUrl)
	if err := service.UpdateProjectInfo(projectName, editedName, avatarUrl, desc); !util.IsOk(err) {
		response.Fail(ctx, nil, "更新项目信息失败！")
		return
	}
	response.Success(ctx, nil, "更新项目信息成功！")
}

// 发送项目列表
func SendProjectTable(ctx *gin.Context) {
	DB = common.GetDB()
	projects, count := service.GetAllProjects()

	// 返回信息
	response.Success(ctx, gin.H{
		"Count":    count,
		"Projects": projects,
	}, "返回项目表单成功！")
}

// 发送一个项目的信息
func SendProjectInfo(ctx *gin.Context) {
	DB = common.GetDB()

	// 获取项目名称
	projectName := ctx.Param("name")
	projectId, err := service.GetProjectIdByName(projectName)
	if !util.IsOk(err) {
		response.Fail(ctx, nil, "项目信息id获取失败！"+err.Message)
		log.Panicln(err)
	}
	// 根据项目名称获取项目
	project, err1 := service.GetProjectByName(projectName)
	if !util.IsOk(err1) {
		response.Fail(ctx, nil, "项目信息获取失败！"+err.Message)
		log.Println("项目信息获取失败！")
		return
	}

	// 获取项目捐赠流水
	transactions, err := service.GetProjectTransactions(projectId)
	if !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		log.Panicln(err)
	}
	if *transactions == nil {
		response.Success(ctx, gin.H{
			"Project":      *project,
			"Transactions": [0]model.TransactionDto{},
		}, "捐赠流水获取成功！")
		return
	}

	response.Success(ctx, gin.H{
		"Project":      *project,
		"Transactions": *transactions,
	}, "项目信息获取成功")
}

// 发送一个用户的所有项目
func SendProjectOfUser(ctx *gin.Context) {
	DB = common.GetDB()
	// 获取用户id
	userId, err := service.GetUserIdByName(ctx.Param("name"))
	if !util.IsOk(err) {
		response.Fail(ctx, nil, "获取用户id失败！")
		log.Panicln(err.Message)
	}

	// 根据用户id获取项目
	projects, count, err1 := service.GetProjectsByUserId(userId)
	if !util.IsOk(err1) {
		response.Fail(ctx, nil, "获取用户的所有项目信息失败！")
		log.Panicln("获取用户的所有项目信息失败！")
	}
	response.Success(ctx, gin.H{
		"Count":    count,
		"Projects": *projects,
	}, "获取用户的所有项目信息成功")
}

// 根据项目id删除项目
func DeleteProject(ctx *gin.Context) {
	DB = common.GetDB()
	projectId, err := service.GetProjectIdByName(ctx.Param("name"))
	if !util.IsOk(err) {
		response.Success(ctx, nil, "项目id获取失败！")
		log.Panicln(err)
	}

	// 删除项目
	if err := service.DeleteProjectById(projectId); !util.IsOk(err) {
		response.Success(ctx, nil, err.Message)
		log.Panicln(err)
	}

	// 删除项目成员
	if err := service.DeleteAllProjectMembers(projectId); !util.IsOk(err) {
		response.Success(ctx, nil, err.Message)
		log.Panicln(err)
	}

	// 删除项目仓库
	if err := service.DeleteProjectRepos(projectId); !util.IsOk(err) {
		response.Success(ctx, nil, err.Message)
		log.Panicln(err)
	}

	// 删除项目钱包
	if err := service.DeleteProjectWallets(projectId); !util.IsOk(err) {
		response.Success(ctx, nil, err.Message)
		log.Panicln(err)
	}

	// 删除成员钱包
	if err := service.DeleteMemberWallets(projectId); !util.IsOk(err) {
		response.Success(ctx, nil, err.Message)
		log.Panicln(err)
	}

	// 删除项目机器人
	if err := service.DeleteProjectBots(projectId); !util.IsOk(err) {
		response.Success(ctx, nil, err.Message)
		log.Panicln(err)
	}

	// 删除项目受捐记录
	if err := service.DeleteTransaction(projectId); !util.IsOk(err) {
		response.Success(ctx, nil, err.Message)
		log.Panicln(err)
	}
	response.Success(ctx, nil, "删除项目成功！")
}
