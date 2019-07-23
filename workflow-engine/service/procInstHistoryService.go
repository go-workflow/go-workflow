package service

import (
	"errors"

	"github.com/go-workflow/go-workflow/workflow-engine/model"
	"github.com/mumushuiding/util"
)

// FindProcHistory 查询我的审批
func FindProcHistory(receiver *ProcessPageReceiver) (string, error) {
	datas, count, err := findAllProcHistory(receiver)
	if err != nil {
		return "", err
	}
	return util.ToPageJSON(datas, count, receiver.PageIndex, receiver.PageSize)
}

// FindProcHistoryByToken 查询我的审批纪录
func FindProcHistoryByToken(token string, receiver *ProcessPageReceiver) (string, error) {
	userinfo, err := GetUserinfoFromRedis(token)
	if err != nil {
		return "", err
	}
	if len(userinfo.Company) == 0 {
		return "", errors.New("保存在redis中的【用户信息 userinfo】字段 company 不能为空")
	}
	if len(userinfo.ID) == 0 {
		return "", errors.New("保存在redis中的【用户信息 userinfo】字段 ID 不能为空")
	}
	receiver.Company = userinfo.Company
	receiver.UserID = userinfo.ID
	// receiver.Username = userinfo.Username
	return FindProcHistory(receiver)
}
func findAllProcHistory(receiver *ProcessPageReceiver) ([]*model.ProcInstHistory, int, error) {
	var page = util.Page{}
	page.PageRequest(receiver.PageIndex, receiver.PageSize)
	return model.FindProcHistory(receiver.UserID, receiver.Company, receiver.PageIndex, receiver.PageSize)
}

// DelProcInstHistoryByID DelProcInstHistoryByID
func DelProcInstHistoryByID(id int) error {
	return model.DelProcInstHistoryByID(id)
}

// StartHistoryByMyself 查询我发起的流程
func StartHistoryByMyself(receiver *ProcessPageReceiver) (string, error) {
	var page = util.Page{}
	page.PageRequest(receiver.PageIndex, receiver.PageSize)
	datas, count, err := model.StartHistoryByMyself(receiver.UserID, receiver.Company, receiver.PageIndex, receiver.PageSize)
	if err != nil {
		return "", err
	}
	return util.ToPageJSON(datas, count, receiver.PageIndex, receiver.PageSize)
}

// FindProcHistoryNotify 查询抄送我的流程
func FindProcHistoryNotify(receiver *ProcessPageReceiver) (string, error) {
	var page = util.Page{}
	page.PageRequest(receiver.PageIndex, receiver.PageSize)
	datas, count, err := model.FindProcHistoryNotify(receiver.UserID, receiver.Company, receiver.Groups, receiver.PageIndex, receiver.PageSize)
	if err != nil {
		return "", err
	}
	return util.ToPageJSON(datas, count, receiver.PageIndex, receiver.PageSize)
}
