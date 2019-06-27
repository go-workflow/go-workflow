package service

import (
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
func findAllProcHistory(receiver *ProcessPageReceiver) ([]*model.ProcInstHistory, int, error) {
	var page = util.Page{}
	page.PageRequest(receiver.PageIndex, receiver.PageSize)
	return model.FindProcHistory(receiver.UserID, receiver.Company, receiver.PageIndex, receiver.PageSize)
}

// DelProcInstHistoryByID DelProcInstHistoryByID
func DelProcInstHistoryByID(id int) error {
	return model.DelProcInstHistoryByID(id)
}
