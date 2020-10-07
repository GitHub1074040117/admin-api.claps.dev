package common

import (
	"claps-admin/model"
	"claps-admin/util"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// 根据此值进行哈希
var jwtKey = []byte("a_secret_crect")

type Claims struct {
	Email string
	jwt.StandardClaims
}

// 创建管理员Token
func ReleaseToken(admin model.Admin) (string, *util.Err) {
	// 1天后过期
	expirationTime := time.Now().Add(1 * 24 * time.Hour)

	//记录在token里面的相关信息
	claims := &Claims{
		Email: admin.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "claps",
			Subject:   "admin_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		log.Panicln("ReleaseToken:Token转化为字符串时发生异常！", err)
	}

	return tokenString, util.Success()

}

// 解析token
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claims, err
}
