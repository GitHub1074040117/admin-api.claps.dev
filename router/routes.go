package router

import (
	"claps-admin/common"
	"claps-admin/controller"
	"claps-admin/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware())
	r.POST("tokenAuth", middleware.TokenAuthHandler)

	r.GET("_hc", common.HealthCheck)

	loginGroup := r.Group("login/")
	{
		loginGroup.GET("authUrl", controller.SendAuthUrl)
		loginGroup.POST("code", controller.CheckLoginCode)
		loginGroup.GET("info", controller.SendLoginInfo)
	}

	r.Use(middleware.AuthMiddleware())
	adminGroup := r.Group("admin/")
	{
		adminGroup.GET("table", controller.SendAdminTable)
		adminGroup.POST("add", controller.AddNewAdmin)
		adminGroup.POST("delete", controller.DeleteAdmin)
	}
	r.POST("addAdmins", controller.AddNewAdmins)

	projectGroup := r.Group("projects")
	{
		projectGroup.POST("/add", controller.AddNewProject)
		projectGroup.GET("/table", controller.SendProjectTable)

		//projectGroup.POST("/repoAdd", controller.AddProjectRepo)
		projectGroup.POST("/repoDelete", controller.DeleteProjectRepository)
	}

	r.GET("project/:name/info", controller.SendProjectInfo)
	r.GET("project/:name/member", controller.SendProjectMembers)
	r.POST("project/:name/edit", controller.UpdateProjectInfo)
	r.GET("project/:name/repo", controller.SendProjectRepos)
	r.POST("project/:name/user", controller.SendProjectOfUser)
	r.POST("project/:name/delete", controller.DeleteProject)
	r.POST("project/:name/memberDelete", controller.DeleteProjectMember)
	r.POST("project/:name/memberAdd", controller.AddProjectMember)
	r.POST("project/:name/addMembers", controller.AddProjectMembers)
	r.POST("project/:name/addRepository", controller.AddProjectRepo)

	r.POST("trans/:name/transactions", controller.SendTransactions)
	r.POST("trans/:name/transfers", controller.SendTransfers)
	r.GET("transactions/all", controller.SendAllTransactions)
	r.GET("transfers/all", controller.SendAllTransfers)

	r.GET("user/:name/info", controller.SendUserInfo)
	r.GET("users/all", controller.SendAllUsers)
	/*r.POST("users/exclude", controller.SendNotProjectUsers)*/
	return r
}
