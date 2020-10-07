package service

import (
	"claps-admin/common"
	"claps-admin/merico"
	"claps-admin/model"
	"claps-admin/util"
	"log"
)

const GithubUrl = "https://github.com/"
const GITHUB = "GITHUB"
const GITLAB = "GITLAB"
const GIT = "GIT"
const BITBUCKET = "BITBUCKET"

// 创建仓库添加到数据库中，返回repoId
func CreateRepository(projectId int64, repoType string, repoUrl string, repoName string) (int64, *util.Err) {
	DB := common.GetDB()
	var newRepo model.Repository
	var ptm model.ProjectToMericoGroup
	var repoToMerico model.RepositoryToMericoProject
	DB.Where("project_id = ?", projectId).First(&ptm)
	if len(ptm.MericoGroupId) == 0 {
		log.Println("在ProjectToMericoGroup表中查找不到该记录：projectId = ", projectId)
		return 0, util.Fail("查找不到该记录")
	}
	// 在merico添加仓库
	merProjectId, err := merico.AddRepository(repoUrl, ptm.MericoGroupId)
	if !util.IsOk(err) {
		return 0, util.Fail("merico添加仓库失败！")
	}
	newRepo.ProjectId = projectId
	newRepo.Type = repoType
	newRepo.Slug = repoUrlCut(repoUrl, repoType)
	newRepo.Name = repoName
	if err := DB.Create(&newRepo).Error; err != nil {
		merico.DeleteRepository(merProjectId)
		return 0, util.Fail(err.Error())
	}
	// 在数据库中添加仓库关系
	repoToMerico.MericoProjectId = merProjectId
	repoToMerico.RepositoryId = newRepo.Id
	if err := DB.Create(&repoToMerico).Error; err != nil {
		merico.DeleteRepository(merProjectId)
		return 0, util.Fail(err.Error())
	}
	return newRepo.Id, util.Success()
}

// 根据项目Id获取所有仓库,返回指针类型
func GetAllReposByProjectId(projectId int64) (*[]model.Repository, *util.Err) {
	DB := common.GetDB()
	// 查仓库表获取某个项目的仓库数组
	var repos []model.Repository
	if err := DB.Where("project_id = ?", projectId).Find(&repos).Error; err != nil {
		return nil, util.Fail("数据库查询repo表失败！" + err.Error())
	}
	if len(repos) == 0 {
		return &repos, util.Fail("获取仓库失败，未查找到任何仓库！")
	}
	// 补充仓库url
	for i := 0; i < len(repos); i++ {
		repos[i].Slug = repoUrlFill(repos[i].Slug, repos[i].Type)
	}
	return &repos, util.Success()
}

// 根据项目id删除项目的所有仓库
func DeleteProjectRepos(projectId int64) *util.Err {
	DB := common.GetDB()
	var repos []model.Repository
	// 在数据库中删除
	DB.Where("project_id = ?", projectId).Delete(&repos)
	// 在merico中删除
	for i := 0; i < len(repos); i++ {
		if err := merico.DeleteRepository(repoIdToMericoId(repos[i].Id)); !util.IsOk(err) {
			log.Println(err.Message)
			return util.Fail(err.Message)
		}
	}
	return util.Success()
}

/*// 判断一个仓库是否在数据库中存在
func IsRepositoryExist(repo *model.Repository) bool {
	var repoInDB model.Repository
	// 按slug查数据库
	DB.Where("slug = ? AND project_id = ? AND type = ?", repo.Slug, repo.ProjectId, repo.Type).First(&repoInDB)
	if len(repoInDB.Id) != 0 {
		return true
	}
	return false
}*/

// 根据仓库id删除仓库
func DeleteRepository(repoId int64) *util.Err {
	DB := common.GetDB()
	var repo model.Repository
	// 在数据库中删除
	if err := DB.Where("id = ?", repoId).Delete(&repo).Error; err != nil {
		log.Println(err)
		return util.Fail(err.Error())
	}
	// 在merico中删除
	if err := merico.DeleteRepository(repoIdToMericoId(repoId)); !util.IsOk(err) {
		log.Println(err.Message)
		return util.Fail(err.Message)
	}
	return util.Success()
}

func repoUrlCut(repoUrl string, repoType string) string {
	switch repoType {
	case GITHUB:
		if len(repoUrl) > 19 {
			return repoUrl[19:]
		}
		break
	default:
		break
	}
	return repoUrl
}

func repoUrlFill(repoUrl string, repoType string) string {
	switch repoType {
	case GITHUB:
		return GithubUrl + repoUrl
	default:
		break
	}
	return repoUrl
}

func repoIdToMericoId(repoId int64) string {
	var repoMerico model.RepositoryToMericoProject
	db := common.GetDB()
	if err := db.Where("repository_id = ?", repoId).First(&repoMerico).Error; err != nil {
		log.Println(err)
		return ""
	}
	return repoMerico.MericoProjectId
}
