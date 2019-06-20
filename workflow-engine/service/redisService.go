package service

import (
	"github.com/mumushuiding/util"

	"github.com/go-workflow/go-workflow/workflow-engine/model"
)

// UserInfo 用户信息
type UserInfo struct {
	Company     string   `json:"company"`
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	Departments []string `json:"departments"`
}

// GetUserinfoFromRedis GetUserinfoFromRedis
func GetUserinfoFromRedis(token string) (*UserInfo, error) {
	result, err := GetValFromRedis(token)
	if err != nil {
		return nil, err
	}
	// fmt.Println(result)
	var userinfo = &UserInfo{}
	err = util.Str2Struct(result, userinfo)
	if err != nil {
		return nil, err
	}
	return userinfo, nil
}

// GetValFromRedis 从redis获取值
func GetValFromRedis(key string) (string, error) {
	return model.RedisGetVal(key)
}
