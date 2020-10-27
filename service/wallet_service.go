package service

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/util"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

const AssetType int = 5

var assets = [AssetType]string{model.BTC, model.BCH, model.ETH, model.XRP, model.DOGE}

// 通过map创建项目钱包，每个币种对应一个项目钱包，因此需要创建5个
func CreateProjectWallets(db *gorm.DB, projectId int64) *util.Err {

	// 将在项目钱包中创建8条记录
	for j := 0; j < AssetType; j++ {
		wallet0, err := createWallet(projectId, model.PersperAlgorithm, assets[j])
		if !util.IsOk(err) {
			log.Panicln("项目钱包创建失败！")
			return err
		}
		db.Create(&wallet0)
		wallet1, err := createWallet(projectId, model.Commits, assets[j])
		if !util.IsOk(err) {
			log.Panicln("项目钱包创建失败！")
			return err
		}
		db.Create(&wallet1)
		wallet2, err := createWallet(projectId, model.ChangedLines, assets[j])
		if !util.IsOk(err) {
			log.Panicln("项目钱包创建失败！")
			return err
		}
		db.Create(&wallet2)
		wallet3, err := createWallet(projectId, model.IdenticalAmount, assets[j])
		if !util.IsOk(err) {
			log.Panicln("项目钱包创建失败！")
			return err
		}
		db.Create(&wallet3)
	}

	return util.Success()
}

// 创建一个项目钱包
func createWallet(projectId int64, distribution string, assetId string) (*model.Wallet, *util.Err) {
	var wallet model.Wallet
	var err *util.Err
	wallet.ProjectId = projectId
	wallet.BotId, err = GetBotId(projectId, distribution)
	if !util.IsOk(err) {
		log.Panicln("BotID 获取失败！")
	}
	wallet.AssetId = assetId
	wallet.SyncedAt = time.Now()
	return &wallet, util.Success()
}

// 根据项目id创建成员钱包,每个成员对应一个钱包，每个币种(5种)对应一条记录
func CreateMemberWallets(db *gorm.DB, projectId int64) *util.Err {
	members, err := GetAllMembersByProjectId(projectId)
	if !util.IsOk(err) {
		return err
	}
	// n个成员生成n*8条记录
	for i := 0; i < len(*members); i++ {
		for j := 0; j < AssetType; j++ {
			if IsMemberWalletExisted(projectId, (*members)[i].Id) {
				continue
			}
			memberWallet0, err := createMemberWallet((*members)[i].Id, projectId, model.PersperAlgorithm, assets[j])
			if !util.IsOk(err) {
				log.Panicln("成员钱包创建失败！", err)
			}
			fmt.Println(memberWallet0.UserId, " ", memberWallet0.ProjectId)
			if err := db.Create(memberWallet0).Error; err != nil {
				log.Panicln(err)
			}
			memberWallet1, err := createMemberWallet((*members)[i].Id, projectId, model.Commits, assets[j])
			if !util.IsOk(err) {
				log.Panicln("成员钱包创建失败！", err)
			}
			if err := db.Create(memberWallet1).Error; err != nil {
				log.Panicln(err)
			}
			memberWallet2, err := createMemberWallet((*members)[i].Id, projectId, model.ChangedLines, assets[j])
			if !util.IsOk(err) {
				log.Panicln("成员钱包创建失败！", err)
			}
			if err := db.Create(memberWallet2).Error; err != nil {
				log.Panicln(err)
			}
			memberWallet3, err := createMemberWallet((*members)[i].Id, projectId, model.IdenticalAmount, assets[j])
			if !util.IsOk(err) {
				log.Panicln("成员钱包创建失败！", err)
			}
			if err := db.Create(memberWallet3).Error; err != nil {
				log.Panicln(err)
			}
		}
	}
	return util.Success()
}

// 为项目添加一个成员钱包
func CreateMemberWallet(projectId int64, userId int64) *util.Err {
	db := common.GetDB()
	distribution := GetDistributionByProjectId(projectId)
	for j := 0; j < AssetType; j++ {
		memberWallet, err := createMemberWallet(userId, projectId, distribution, assets[j])
		if !util.IsOk(err) {
			log.Panicln("成员钱包创建失败！", err)
		}
		db.Create(memberWallet)
	}
	return util.Success()
}

// 创建一个成员钱包
func createMemberWallet(userId int64, projectId int64, distribution string, assetId string) (*model.MemberWallet, *util.Err) {
	var memberWallet model.MemberWallet
	var err *util.Err
	memberWallet.UserId = userId
	memberWallet.ProjectId = projectId
	memberWallet.AssetId = assetId
	memberWallet.BotId, err = GetBotId(projectId, distribution)
	if !util.IsOk(err) {
		log.Panicln("BotID 获取失败！")
	}
	return &memberWallet, util.Success()
}

// 根据项目id删除所有项目钱包
func DeleteProjectWallets(projectId int64) *util.Err {
	db := common.GetDB()
	var wallets []model.Wallet
	db.Where("project_id = ?", projectId).Delete(&wallets)
	return util.Success()
}

// 根据项目id删除所有成员钱包
func DeleteMemberWallets(projectId int64) *util.Err {
	db := common.GetDB()
	var memberWallets []model.MemberWallet
	db.Where("project_id = ?", projectId).Delete(&memberWallets)
	return util.Success()
}

// 根据项目id获取分配方式类型
func GetDistributionByProjectId(projectId int64) string {
	db := common.GetDB()
	var bot model.Bot
	err := db.Where("project_id = ?", projectId).First(&bot).Error
	if err != nil {
		log.Panicln(err)
	}
	return bot.Distribution
}

// 根据项目id和用户id查找成员钱包是否存在
func IsMemberWalletExisted(projectId int64, userId int64) bool {
	db := common.GetDB()
	var mw model.MemberWallet
	db.Where("project_id = ? AND user_id = ?", projectId, userId).First(&mw)
	if mw.ProjectId == 0 {
		return false
	}
	return true
}
