package main

import (
	"claps-admin/common"
	"claps-admin/router"
	"claps-admin/service"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"os"
)

func main() {

	//初始化读取配置
	InitConfig()
	db := common.InitDB()
	defer db.Close()

	service.InitAdmin()
	r := gin.Default()
	r = router.CollectRoute(r)
	port := viper.GetString("server.port")
	if port != "" {
		//取代默认端口
		panic(r.Run(":" + port))
	}
	panic(r.Run())

}

// 数据库的配置文件
func InitConfig() {
	workDir, _ := os.Getwd()
	//设置要读取的文件名
	viper.SetConfigName("application")
	//设置读取文件的类型
	viper.SetConfigType("yml")
	//设置文件的路径
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("读取配置文件失败")
	}
}
