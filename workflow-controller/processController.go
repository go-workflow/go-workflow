package controller

import (
	"fmt"
	"net/http"

	"github.com/mumushuiding/util"

	"github.com/go-workflow/go-workflow/workflow-engine/service"
)

// StartProcessInstance 启动流程
func StartProcessInstance(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only suppoert Post ")
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
	if len(proc.Title) == 0 {
		util.Response(writer, "启动流程的标题title不能为空", false)
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
	return
}

// FindMyProcInstPageAsJSON FindMyProcInstPageAsJSON
// 查询到我审批的流程实例
func FindMyProcInstPageAsJSON(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only suppoert Post ")
		return
	}
	var receiver = service.GetDefaultProcessPageReceiver()
	err := util.Body2Struct(request, &receiver)
	if err != nil {
		util.ResponseErr(writer, err)
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
	util.Response(writer, result, true)
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
