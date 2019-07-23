package service

import (
	"github.com/mumushuiding/util"

	"github.com/go-workflow/go-workflow/workflow-engine/model"
)

// UserInfo 用户信息
type UserInfo struct {
	Company string `json:"company"`
	// 用户所属部门
	Department string `json:"department"`
	Username   string `json:"username"`
	ID         string `json:"ID"`
	// 用户的角色
	Roles []string `json:"roles"`
	// 用户负责的部门
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
