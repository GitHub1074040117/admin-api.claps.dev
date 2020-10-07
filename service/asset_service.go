package service

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/util"
	"log"
)

func GetAssetById(assetId string) (*model.Asset, *util.Err) {
	db := common.GetDB()
	var asset model.Asset
	err := db.Where("asset_id = ?", assetId).First(&asset).Error
	if err != nil {
		log.Println(err)
		return nil, util.Fail(err.Error())
	}
	if len(asset.AssetId) == 0 {
		log.Println("数据库查找不到该币种: ", assetId)
		return nil, util.Fail("数据库查找不到该币种")
	}
	return &asset, util.Success()
}
