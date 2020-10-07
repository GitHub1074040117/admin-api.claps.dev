package service

import (
	"claps-admin/model"
	"claps-admin/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// 通过code获取token认证url
func GetTokenAuthUrl(ClientId string, ClientSecret string, code string) string {
	return fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		ClientId, ClientSecret, code,
	)
}

// 获取 token
func GetToken(url string) (*oauth2.Token, *util.Err) {

	// 形成请求
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, url, nil); err != nil {
		log.Println(err)
		return nil, util.Fail("请求形成失败")
	}
	req.Header.Set("accept", "application/json")

	// 发送请求并获得响应
	var httpClient = http.Client{}
	var res *http.Response
	if res, err = httpClient.Do(req); err != nil {
		log.Println(err)
		return nil, util.Fail("请求发送失败")
	}

	// 将响应体解析为 token，并返回
	var token oauth2.Token
	if err = json.NewDecoder(res.Body).Decode(&token); err != nil {
		log.Println(err)
		return nil, util.Fail("token解析失败")
	}
	//fmt.Printf("%+v\n", token)

	return &token, util.Success()
}

func GetUserInfo(token *oauth2.Token) (*github.User, *util.Err) {

	/*var userInfoUrl = "https://api.github.com/user"
	var req *http.Request
	var err error
	if req, err = http.NewRequest(http.MethodGet, userInfoUrl, nil); err != nil {
		log.Panicln(err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token.AccessToken))

	var client = http.Client{}
	var res *http.Response
	if res, err = client.Do(req); err != nil {
		log.Panicln(err)
	}

	var userInfo = make(map[string]interface{})
	if err = json.NewDecoder(res.Body).Decode(&userInfo); err != nil {
		fmt.Println(userInfo)
		log.Panicln(err)
	}
	fmt.Println(userInfo)

	var user github.User
	return &user, util.Success()*/

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token.AccessToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")

	if err != nil {
		log.Println(err)
		return user, util.Fail("获取用户信息失败!")
	}

	return user, util.Success()
}

// 更新用户在数据库中的信息
func UpdateUserInfo(DB *gorm.DB, user model.Admin) {
	DB.Model(&user).Where("email = ?", user.Email).Updates(user)
}
