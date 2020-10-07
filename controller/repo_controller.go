package controller

import (
	"claps-admin/common"
	"claps-admin/response"
	"claps-admin/service"
	"claps-admin/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

// 从数据库读取某个项目的所有仓库，并发送给前端
func SendProjectRepos(ctx *gin.Context) {
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

	// 获取项目仓库数组
	repos, err := service.GetAllReposByProjectId(projectId)
	if !util.IsOk(err) {
		response.Fail(ctx, nil, err.Message)
		log.Panicln(err)
	}

	// 发送数据
	count := len(*repos)
	response.Success(ctx, gin.H{
		"Count": count,
		"Repos": *repos,
	}, "获取仓库组成功！")
}

// 为项目添加仓库
func AddProjectRepo(ctx *gin.Context) {
	// 解析
	var infoMap = make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&infoMap)
	repoName := infoMap["repoName"]
	//description := infoMap["description"]
	repoType := infoMap["repoType"]
	repoUrl := infoMap["repoUrl"]
	projectName := ctx.Param("name")
	projectId, err := service.GetProjectIdByName(projectName)
	if !util.IsOk(err) {
		response.Fail(ctx, nil, "项目id获取失败！")
		log.Panicln(err)
	}
	if _, err = service.CreateRepository(projectId, repoType, repoUrl, repoName); !util.IsOk(err) {
		response.Fail(ctx, nil, "仓库添加失败！")
		log.Panicln(err)
	}
	response.Success(ctx, nil, "仓库添加成功！")
}

// 根据仓库id删除仓库
func DeleteProjectRepository(ctx *gin.Context) {
	var infoMap = make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&infoMap)
	// 解析获取仓库id
	repoId, err := strconv.Atoi(infoMap["repoId"])
	if err != nil {
		response.Fail(ctx, nil, "项目删除失败！")
		log.Panicln(err)
	}
	// 删除仓库
	err1 := service.DeleteRepository(int64(repoId))
	if !util.IsOk(err1) {
		response.Fail(ctx, nil, "项目删除失败！")
		log.Panicln(err1)
	}
	response.Success(ctx, nil, "删除成功！")
}
