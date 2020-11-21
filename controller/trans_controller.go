package controller

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/response"
	"claps-admin/service"
	"claps-admin/util"
	"github.com/gin-gonic/gin"
	"log"
)

// 发送某个项目的受捐流水
func SendTransactions(ctx *gin.Context) {
	DB = common.GetDB()
	// 获取传送过来的项目id
	projectId, _ := service.GetProjectIdByName(ctx.Param("name"))

	// 获取受捐记录
	transactions, count := service.GetProjectTransactions(projectId)
	response.Success(ctx, gin.H{
		"Count":        count,
		"Transactions": transactions,
	}, "返回项目受捐流水成功！")
}

// 发送所有项目的受捐流水
func SendAllTransactions(ctx *gin.Context) {
	transactions, err := service.GetAllTransactions()
	if !util.IsOk(err) {
		response.Fail(ctx, nil, "受捐流水获取失败！")
		log.Panicln(err.Message)
	}
	if *transactions == nil {
		response.Success(ctx, gin.H{
			"Transactions": [0]model.TransactionDto{},
		}, "捐赠流水获取成功！")
	} else {
		response.Success(ctx, gin.H{
			"Transactions": transactions,
		}, "捐赠流水获取成功！")
	}

}

// 发送所有用户的提现流水
func SendAllTransfers(ctx *gin.Context) {
	transfers, err := service.GetAllTransfers()
	if !util.IsOk(err) {
		response.Fail(ctx, nil, "获取提现流水失败！")
		log.Panicln(err.Message)
	}
	if *transfers == nil {
		response.Success(ctx, gin.H{
			"Transfer": [0]model.TransferDto{},
		}, "提现流水获取成功！")
	} else {
		response.Success(ctx, gin.H{
			"Transfers": transfers,
		}, "提现流水获取成功！")
	}
}

// 发送某个用户的提现记录
func SendTransfers(ctx *gin.Context) {
	DB = common.GetDB()
	userId, _ := service.GetUserIdByName(ctx.Param("name"))
	user, err := service.GetUserById(userId)
	if !util.IsOk(err) {
		response.Fail(ctx, nil, "获取用户信息失败！")
		return
	}

	// 获取提现记录
	transfers, count, err1 := service.GetUserTransfers(user.MixinId)
	if !util.IsOk(err1) {
		response.Fail(ctx, nil, "获取用户提现记录失败！")
		return
	}
	response.Success(ctx, gin.H{
		"Count":     count,
		"Transfers": transfers,
	}, "返回项目提现记录成功！")
}
