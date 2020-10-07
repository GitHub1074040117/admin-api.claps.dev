package middleware

import (
	"claps-admin/common"
	"claps-admin/model"
	"claps-admin/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

const AuthOff = "off"

// gin的中间件，返回一个handler, 实现token的认证和保护路由
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		//fmt.Println("token: "+tokenString)
		// validate token format
		if len(tokenString) < 8 {
			response.Fail(ctx, nil, "Invalid redirect!")
			log.Println("token为空!非法访问！: " + tokenString)
			return
		}
		// 从第7位开始向后截取,解析token
		//tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil {
			response.Fail(ctx, nil, "Invalid token")
			log.Println("解析token错误:", err)
			return
		}
		if !token.Valid {
			response.Fail(ctx, nil, "Invalid redirect!")
			log.Println("token无效!非法访问！")
			return
		}

		// 若不需要验证身份则通过
		if viper.GetString("authOff") == AuthOff {
			response.Success(ctx, nil, "")
			ctx.Next()
			return
		}

		// 验证通过后获取claim中的userId
		email := claims.Email
		DB := common.GetDB()
		var admin model.Admin
		DB.Where("email = ?", email).First(&admin)
		//DB.First(&admin, email)

		// 如果用户不存在,返回
		if admin.UserId == 0 {
			response.Fail(ctx, nil, "Invalid token")
			log.Println("用户不存在!非法访问！")
		}

		// 用户存在，将User信息写入上下文
		ctx.Set("user", admin)
		ctx.Next()
	}
}

// 验证token Handler
func TokenAuthHandler(ctx *gin.Context) {
	// 获取token
	var requestMap = make(map[string]string)
	json.NewDecoder(ctx.Request.Body).Decode(&requestMap)
	var tokenString = requestMap["token"]

	// validate token format
	if len(tokenString) < 8 {
		response.Fail(ctx, nil, "Invalid redirect!")
		log.Println("token为空!非法访问！: " + tokenString)
		return
	}
	// 从第7位开始向后截取,解析token
	//tokenString = tokenString[7:]
	token, claims, err := common.ParseToken(tokenString)
	if err != nil {
		response.Fail(ctx, nil, "Invalid token")
		log.Println(err)
		return
	}
	if !token.Valid {
		response.Fail(ctx, nil, "Invalid redirect!")
		log.Println("token无效!非法访问！")
		return
	}

	// 若不需要验证身份则通过
	if viper.GetString("authOff") == AuthOff {
		response.Success(ctx, nil, "")
		ctx.Next()
		return
	}

	// 验证通过后获取claim中的userId
	email := claims.Email
	DB := common.GetDB()
	var admin model.Admin
	DB.Where("email = ?", email).First(&admin)
	//DB.First(&admin, email)

	// 如果用户不存在,返回
	if admin.UserId == 0 {
		response.Fail(ctx, nil, "Invalid token")
		log.Println("用户不存在!非法访问！")
	}

	response.Success(ctx, nil, "")
	// 用户存在，将User信息写入上下文
	ctx.Set("user", admin)
	ctx.Next()
}
