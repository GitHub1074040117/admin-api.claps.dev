package service

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/util"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"flag"
	"github.com/fox-one/mixin-sdk-go"
	"log"
	"os"
)

var config = flag.String("config", "./config/keystore-7000103520.json", "keystore file path")

const DistributionType = 4

var distributionSet = [DistributionType]string{model.PersperAlgorithm, model.Commits, model.ChangedLines, model.IdenticalAmount}

// 为项目创建四个机器人
func Create4BotsForProject(projectId int64) *util.Err {

	db := common.GetDB()
	flag.Parse()
	// 4种不同的分配算法
	for i := 0; i < DistributionType; i++ {
		bot, err := createBot(config, projectId, distributionSet[i])
		if !util.IsOk(err) {
			return err
		}
		db.Create(bot)
	}
	return util.Success()
}

// 为项目添加一个机器人， 并设置其分配类型
func createBot(config *string, projectId int64, distribution string) (*model.Bot, *util.Err) {
	ctx := context.Background()
	f, err := os.Open(*config)
	if err != nil {
		log.Println("bot_service/createBot: 配置文件打开错误: ", err)
		return nil, util.Fail("bot_service/createBot: 配置文件打开错误")
	}

	var store mixin.Keystore
	if err := json.NewDecoder(f).Decode(&store); err != nil {
		log.Println("bot_service/createBot: keystore解析错误: ", err)
		return nil, util.Fail("bot_service/createBot: keystore解析错误")
	}

	// 根据keystore创建client主机器人
	client, err := mixin.NewFromKeystore(&store)
	if err != nil {
		log.Println("bot_service/createBot: 创建client主机器人错误: ", err)
		return nil, util.Fail("bot_service/createBot: 创建client主机器人错误")
	}

	// 由主机器人创建子机器人的keystore
	privateKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	_, subStore, err := client.CreateUser(ctx, privateKey, "sub robot")
	if err != nil {
		log.Println("bot_service/createBot: 创建子机器人的keystore错误: ", err)
		return nil, util.Fail("bot_service/createBot: 创建子机器人的keystore错误")
	}

	// 由子机器人的keystore创建子机器人，并设置pin
	newPin := mixin.RandomPin()
	subClient, _ := mixin.NewFromKeystore(subStore)
	if err := subClient.ModifyPin(ctx, "", newPin); err != nil {
		log.Println("bot_service/createBot: 创建子机器人的pin错误: ", err)
		return nil, util.Fail("bot_service/createBot: 创建子机器人的pin错误")
	}

	// 验证创建子机器人的pin的有效性
	if err := subClient.VerifyPin(ctx, newPin); err != nil {
		log.Println("子机器人的pin的不具备有效性: ", err)
	}

	var bot model.Bot
	bot.Pin = newPin
	bot.Id = subStore.ClientID
	bot.SessionId = subStore.SessionID
	bot.PinToken = subStore.PinToken
	bot.PrivateKey = subStore.PrivateKey
	bot.Distribution = distribution
	bot.ProjectId = projectId
	return &bot, util.Success()
}

// 根据项目id删除所有机器人
func DeleteProjectBots(projectId int64) *util.Err {
	db := common.GetDB()
	var bots []model.Bot
	db.Where("project_id = ?", projectId).Delete(&bots)
	return util.Success()
}

// 根据项目id和分配方式获取机器人id
func GetBotId(projectId int64, distribution string) (string, *util.Err) {
	var bot model.Bot
	db := common.GetDB()
	if distribution == "PersperAlgorithm" {
		distribution = model.PersperAlgorithm
	}
	if distribution == "Commits" {
		distribution = model.Commits
	}
	if distribution == "ChangedLines" {
		distribution = model.ChangedLines
	}
	if distribution == "IdenticalAmount" {
		distribution = model.IdenticalAmount
	}
	err := db.Where("project_id = ? AND distribution = ?", projectId, distribution).First(&bot).Error
	if bot.Id == "" {
		log.Panicln("未能查找到bot", projectId, distribution)
	}
	if err != nil {
		log.Panicln(err)
	}
	return bot.Id, util.Success()
}
