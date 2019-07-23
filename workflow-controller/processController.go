package controller

import (
	"fmt"
	"net/http"

	"github.com/go-workflow/go-workflow/workflow-engine/model"

	"github.com/mumushuiding/util"

	"github.com/go-workflow/go-workflow/workflow-engine/service"
)

// StartProcessInstanceByToken 启动流程
func StartProcessInstanceByToken(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only suppoert Post ")
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
	var proc = service.ProcessReceiver{}
	err := util.Body2Struct(request, &proc)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	if len(proc.ProcName) == 0 {
		util.Response(writer, "流程定义名procName不能为空", false)
		return
	}
	id, err := service.StartProcessInstanceByToken(token, &proc)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.Response(writer, fmt.Sprintf("%d", id), true)
}

// StartProcessInstance 启动流程
func StartProcessInstance(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only suppoert Post ")
		return
	}
	if model.RedisOpen {
		util.ResponseErr(writer, "已经连接 redis，请使用/workflow/process/startByToken 路径访问")
		return
	}
	var proc = service.ProcessReceiver{}
	err := util.Body2Struct(request, &proc)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	if len(proc.ProcName) == 0 {
		util.Response(writer, "流程定义名procName不能为空", false)
		return
	}
	if len(proc.Company) == 0 {
		util.Response(writer, "用户所在的公司company不能为空", false)
		return
	}
	if len(proc.UserID) == 0 {
		util.Response(writer, "启动流程的用户userId不能为空", false)
		return
	}
	if len(proc.Department) == 0 {
		util.Response(writer, "用户所在部门department不能为空", false)
		return
	}
	id, err := proc.StartProcessInstanceByID(proc.Var)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.Response(writer, fmt.Sprintf("%d", id), true)
}

// FindMyProcInstPageAsJSON FindMyProcInstPageAsJSON
// 查询到我审批的流程实例
func FindMyProcInstPageAsJSON(writer http.ResponseWriter, request *http.Request) {
	if model.RedisOpen {
		util.ResponseErr(writer, "已经连接 redis，请使用/workflow/process/findTaskByToken 路径访问")
		return
	}
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
	result, err := service.FindAllPageAsJSON(receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	fmt.Fprintf(writer, result)
}

// FindMyProcInstByToken FindMyProcInstByToken
// 查询待办的流程
func FindMyProcInstByToken(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！")
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
	// fmt.Printf("token:%s\n", token)
	var receiver = service.GetDefaultProcessPageReceiver()
	err := util.Body2Struct(request, &receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	result, err := service.FindMyProcInstByToken(token, receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	fmt.Fprintf(writer, result)
}

// StartByMyself 我启动的流程
func StartByMyself(writer http.ResponseWriter, request *http.Request) {
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
	result, err := service.StartByMyself(receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	fmt.Fprintf(writer, result)
}

// FindProcNotify 查询抄送我的流程
func FindProcNotify(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持POST方法")
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
	result, err := service.FindProcNotify(receiver)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	fmt.Fprintf(writer, result)
}

// MoveFinishedProcInstToHistory MoveFinishedProcInstToHistory
// 将已经结束的流程实例移动到历史数据库
func MoveFinishedProcInstToHistory(writer http.ResponseWriter, request *http.Request) {
	err := service.MoveFinishedProcInstToHistory()
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseOk(writer)
}
