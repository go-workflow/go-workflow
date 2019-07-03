package controller

import (
	"fmt"
	"net/http"

	"github.com/go-workflow/go-workflow/workflow-engine/model"

	"github.com/go-workflow/go-workflow/workflow-engine/service"
	"github.com/mumushuiding/util"
)

// FindProcHistoryByToken 查看我审批的纪录
func FindProcHistoryByToken(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持POST方法")
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
	var receiver = service.GetDefaultProcessPageReceiver()
	err := util.Body2Struct(request, &receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	result, err := service.FindProcHistoryByToken(token, receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	fmt.Fprintf(writer, result)
}

// FindProcHistory 查询我的审批纪录
func FindProcHistory(writer http.ResponseWriter, request *http.Request) {
	if model.RedisClient != nil {
		util.ResponseErr(writer, "已经连接 redis，请使用/workflow/procHistory/findTaskByToken")
		return
	}
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持POST方法")
		return
	}
	var receiver = service.GetDefaultProcessPageReceiver()
	err := util.Body2Struct(request, &receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	if len(receiver.UserID) == 0 {
		util.Response(writer, "用户userID不能为空", false)
		return
	}
	if len(receiver.Company) == 0 {
		util.Response(writer, "字段 company 不能为空", false)
		return
	}
	result, err := service.FindProcHistory(receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	fmt.Fprintf(writer, result)
}

// StartHistoryByMyself 查询我发起的流程
func StartHistoryByMyself(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only suppoert Post ")
		return
	}
	var receiver = service.GetDefaultProcessPageReceiver()
	err := util.Body2Struct(request, &receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	if len(receiver.UserID) == 0 {
		util.Response(writer, "用户userID不能为空", false)
		return
	}
	if len(receiver.Company) == 0 {
		util.Response(writer, "字段 company 不能为空", false)
		return
	}
	result, err := service.StartHistoryByMyself(receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	fmt.Fprintf(writer, result)
}
