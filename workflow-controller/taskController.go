package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-workflow/go-workflow/workflow-engine/service"

	"github.com/mumushuiding/util"
)

// WithDrawTask 撤回
func WithDrawTask(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only support Post ")
		return
	}
	var taskRe = service.TaskReceiver{}
	err := util.Body2Struct(request, &taskRe)
	str, _ := util.ToJSONStr(taskRe)
	log.Println(str)
	if taskRe.TaskID == 0 {
		util.ResponseErr(writer, "字段taskID不能为空,必须为数字！")
		return
	}
	if len(taskRe.UserID) == 0 {
		util.ResponseErr(writer, "字段userID不能为空！")
		return
	}
	if taskRe.ProcInstID == 0 {
		util.ResponseErr(writer, "字段 procInstID 不能为空,必须为数字！")
		return
	}
	if len(taskRe.Company) == 0 {
		util.ResponseErr(writer, "字段company不能为空！")
		return
	}
	err = service.WithDrawTask(taskRe.TaskID, taskRe.ProcInstID, taskRe.UserID, taskRe.Company)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseOk(writer)
}

// CompleteTask CompleteTask
// 审批
func CompleteTask(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only support Post ")
		return
	}
	var taskRe = service.TaskReceiver{}
	err := util.Body2Struct(request, &taskRe)
	// str, _ := util.ToJSONStr(taskRe)
	// log.Println(str)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	if len(taskRe.Pass) == 0 {
		util.ResponseErr(writer, "字段pass不能为空！")
		return
	}
	pass, err := strconv.ParseBool(taskRe.Pass)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	if taskRe.TaskID == 0 {
		util.ResponseErr(writer, "字段taskID不能为空！")
		return
	}
	if len(taskRe.UserID) == 0 {
		util.ResponseErr(writer, "字段userID不能为空！")
		return
	}
	if len(taskRe.Company) == 0 {
		util.ResponseErr(writer, "字段company不能为空！")
		return
	}
	err = service.Complete(taskRe.TaskID, taskRe.UserID, taskRe.Company, pass)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseOk(writer)
}
