package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-workflow/go-workflow/workflow-engine/model"

	"github.com/go-workflow/go-workflow/workflow-engine/service"

	"github.com/mumushuiding/util"
)

// WithDrawTask 撤回
func WithDrawTask(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only support Post ")
		return
	}
	if model.RedisOpen {
		util.ResponseErr(writer, "已经连接redis缓存，请使用方法 /workflow/task/withdrawByToken ")
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
	err = service.WithDrawTask(taskRe.TaskID, taskRe.ProcInstID, taskRe.UserID, taskRe.Company, taskRe.Comment)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseOk(writer)
}

// WithDrawTaskByToken 撤回
func WithDrawTaskByToken(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only support Post ")
		return
	}
	token := request.Header.Get("Authorization")
	if len(token) == 0 {
		request.ParseForm()
		if len(request.Form["token"]) == 0 {
			util.ResponseErr(writer, "header Authorization 没有保存 token, url参数也不存在 token， 访问失败 ！")
			return
		}
		token = request.Form["token"][0]
	}
	var taskRe = service.TaskReceiver{}
	err := util.Body2Struct(request, &taskRe)
	if taskRe.TaskID == 0 {
		util.ResponseErr(writer, "字段taskID不能为空,必须为数字！")
		return
	}
	if taskRe.ProcInstID == 0 {
		util.ResponseErr(writer, "字段 procInstID 不能为空,必须为数字！")
		return
	}
	err = service.WithDrawTaskByToken(token, &taskRe)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseOk(writer)
}

// CompleteTaskByToken 使用redis缓存时使用当前方法，更安全
func CompleteTaskByToken(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only support Post ")
		return
	}
	token := request.Header.Get("Authorization")
	if len(token) == 0 {
		request.ParseForm()
		if len(request.Form["token"]) == 0 {
			util.ResponseErr(writer, "header Authorization 没有保存 token, url参数也不存在 token， 访问失败 ！")
			return
		}
		token = request.Form["token"][0]
	}
	var taskRe = service.TaskReceiver{}
	err := util.Body2Struct(request, &taskRe)
	// str, _ := util.ToJSONStr(taskRe)
	// log.Println(str)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	if len(taskRe.Comment) > 255 {
		util.ResponseErr(writer, "字段comment 长度不能超过255")
		return
	}
	if len(taskRe.Pass) == 0 {
		util.ResponseErr(writer, "字段pass不能为空！")
		return
	}
	if taskRe.TaskID == 0 {
		util.ResponseErr(writer, "字段taskID不能为空！")
		return
	}
	err = service.CompleteByToken(token, &taskRe)
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
	if model.RedisOpen {
		util.ResponseErr(writer, "已经连接redis缓存，请使用方法 /workflow/task/completeByToken")
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
	err = service.Complete(taskRe.TaskID, taskRe.UserID, taskRe.Company, taskRe.Comment, pass)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseOk(writer)
}
