package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 解决浏览器跨域问题

func CORSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 填写可以访问的域名， *代表所有域名都可以访问
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		// 设置缓存时间
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		// 设置可以通过访问的方法
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// 判断请求是否为options请求
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}

	}
}
