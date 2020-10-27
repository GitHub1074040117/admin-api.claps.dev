package service

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/util"
	"log"
)

// 从数据库中获取某个项目的捐赠流水
func GetProjectTransactions(projectId int64) (*[]model.TransactionDto, *util.Err) {
	db := common.GetDB()
	var transactions []model.Transaction
	var tranDtos []model.TransactionDto
	if err := db.Where("project_id = ?", projectId).Find(&transactions).Error; err != nil {
		log.Println(err)
		return &tranDtos, util.Fail("查询数据库出错！")
	}
	for _, tran := range transactions {
		var tranDto *model.TransactionDto
		tranDto = toTransactionDto(&tran)
		tranDtos = append(tranDtos, *tranDto)
	}
	return &tranDtos, util.Success()
}

// 获取所有项目的捐赠流水
func GetAllTransactions() (*[]model.TransactionDto, *util.Err) {
	var tranDtos []model.TransactionDto
	var transactions []model.Transaction
	DB := common.GetDB()
	if err := DB.Find(&transactions).Error; err != nil {
		log.Println(err)
		return &tranDtos, util.Fail(err.Error())
	}
	if len(transactions) == 0 {
		log.Println("捐赠流水为空")
		return &tranDtos, util.Success()
	}
	for _, tran := range transactions {
		var tranDto *model.TransactionDto
		tranDto = toTransactionDto(&tran)
		tranDtos = append(tranDtos, *tranDto)
	}
	return &tranDtos, util.Success()
}

// 获取所有用户提现记录
func GetAllTransfers() (*[]model.TransferDto, *util.Err) {
	var transfers []model.Transfer
	var tranDtos []model.TransferDto
	DB := common.GetDB()
	if err := DB.Find(&transfers).Error; err != nil {
		log.Panicln(err)
	}
	if len(transfers) == 0 {
		log.Println("提现记录为空")
		return &tranDtos, util.Success()
	}
	for _, tran := range transfers {
		var tranDto *model.TransferDto
		tranDto = ToTransferDto(&tran)
		tranDtos = append(tranDtos, *tranDto)
	}
	return &tranDtos, util.Success()
}

// 从数据库中获取某个用户的提现记录
func GetUserTransfers(mixinId string) (*[]model.Transfer, int, *util.Err) {
	db := common.GetDB()
	var transfers []model.Transfer
	db.Where("mixin_id = ?", mixinId).Find(&transfers)
	count := len(transfers)
	if count == 0 {
		log.Panicln("未查找到该用户的提现记录")
	}
	return &transfers, count, util.Success()
}

// 将transaction转化成transactionDto
func toTransactionDto(transaction *model.Transaction) *model.TransactionDto {
	var transactionDto model.TransactionDto
	project, err := GetProjectById(transaction.ProjectId)
	if !util.IsOk(err) {
		log.Println(err)
		return nil
	}
	user, err := GetUserByMixinId(transaction.Sender)
	if !util.IsOk(err) {
		log.Println("用户不存在！", transaction.Sender)
		transactionDto.Sender = "Not Exist"
	} else {
		transactionDto.SenderAvatar = user.AvatarUrl
		transactionDto.Sender = user.Name
	}
	asset, err := GetAssetById(transaction.AssetId)
	if !util.IsOk(err) {
		return nil
	}

	/*receiver, err := GetProjectByMericoId(transaction.Receiver)
	if !util.IsOk(err) {
		log.Println("由mericoId获取项目失败！")
	}
	fmt.Println(receiver)*/

	transactionDto.Id = transaction.Id
	transactionDto.ProjectAvatar = project.AvatarUrl
	transactionDto.ProjectName = project.Name
	transactionDto.AssetAvatar = asset.IconUrl
	transactionDto.Asset = asset.Name
	transactionDto.Amount = transaction.Amount
	transactionDto.CreatedAt = transaction.CreatedAt
	transactionDto.Receiver = transaction.Receiver

	return &transactionDto
}

// 将transfer转化成transferDto
func ToTransferDto(transfer *model.Transfer) *model.TransferDto {
	var transferDto model.TransferDto
	user, err := GetUserByMixinId(transfer.MixinId)
	if !util.IsOk(err) {
		log.Println(err)
		return &transferDto
	}
	asset, err := GetAssetById(transfer.AssetId)
	if !util.IsOk(err) {
		log.Println("获取assetID出错！")
		return &transferDto
	}
	transferDto.User = user.Name
	transferDto.Asset = asset.Name
	transferDto.Amount = transfer.Amount
	transferDto.CreatedAt = transfer.CreatedAt
	transferDto.Memo = transfer.Memo
	transferDto.SnapshotId = transfer.SnapshotId
	transferDto.TraceId = transfer.TraceId
	return &transferDto
}
