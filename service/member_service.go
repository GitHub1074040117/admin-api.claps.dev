package service

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/util"
	"log"
)

// 根据map中的projectName为指定项目创建成员
func CreateMembers(members interface{}, projectName string) *util.Err {

	DB := common.GetDB()
	var err *util.Err
	for _, userName := range members.([]interface{}) {
		var member model.Member
		member.UserId, err = GetUserIdByName(userName.(string))
		if !util.IsOk(err) {
			log.Panicln("用户id查找失败!")
		}
		member.ProjectId, err = GetProjectIdByName(projectName)
		if !util.IsOk(err) {
			log.Panicln("项目id查找失败!")
		}
		if !IsMemberExist(member.ProjectId, member.UserId) {
			DB.Create(&member)
		}
	}
	return util.Success()
}

// 根据项目Id获取所有成员,返回指针类型
func GetAllMembersByProjectId(projectId int64) (*[]model.User, *util.Err) {
	// 查members表获取映射关系
	db := common.GetDB()
	var members []model.Member
	db.Where("project_id = ?", projectId).Find(&members)
	count := len(members)

	// 查users表获取成员详细信息
	var users []model.User
	for i := 0; i < count; i++ {
		userId := members[i].UserId

		// 数据库中查找用户
		user, err := GetUserById(userId)
		if !util.IsOk(err) {
			return &users, err
		}

		// 在切片中追加用户
		users = append(users, *user)
	}
	return &users, util.Success()
}

// 根据项目id删除所有成员
func DeleteAllProjectMembers(projectId int64) *util.Err {
	db := common.GetDB()
	var members []model.Member
	db.Where("project_id = ?", projectId).Delete(&members)
	return util.Success()
}

/*// 获取传送过来的成员email
func GetMemberEmail(ctx *gin.Context) (string, *util.Err) {
	var emailMap = make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&emailMap)
	memberEmail := emailMap["memberEmail"]
	if len(memberEmail) == 0 {
		fmt.Println("member_service/GetMemberEmail:成员邮箱获取失败！")
		return memberEmail, util.Fail("成员邮箱获取失败！", nil)
	}
	return memberEmail, util.Success()
}*/

// 根据用户id和项目id删除成员
func DeleteMember(userId int64, projectId int64) *util.Err {
	db := common.GetDB()
	// 在members表中删除
	var member model.Member
	db.Where("project_id = ? AND user_id = ?", projectId, userId).Delete(&member)

	// 在member_wallet表中删除
	var memberWallet model.MemberWallet
	db.Where("project_id = ? AND user_id = ?", projectId, userId).Delete(&memberWallet)

	return util.Success()
}

// 判断成员是否存在
func IsMemberExist(projectId int64, userId int64) bool {
	var member model.Member
	db := common.GetDB()
	db.Where("project_id = ? AND user_id = ?", projectId, userId).First(&member)
	if member.ProjectId == 0 {
		return false
	}
	return true
}
