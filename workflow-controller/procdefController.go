package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-workflow/go-workflow/workflow-engine/service"

	"github.com/mumushuiding/util"
)

// SaveProcdef save new procdefnition
// 保存流程定义
func SaveProcdef(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		util.ResponseErr(writer, "只支持Post方法！！Only support Post ")
		return
	}
	var procdef = service.Procdef{}
	err := util.Body2Struct(request, &procdef)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	if len(procdef.Userid) == 0 {
		util.ResponseErr(writer, "字段 userid 不能为空")
		return
	}
	if len(procdef.Company) == 0 {
		util.ResponseErr(writer, "字段 company 不能为空")
		return
	}
	if procdef.Resource == nil {
		util.ResponseErr(writer, "字段 resource 不能为空")
		return
	}
	id, err := procdef.SaveProcdef()
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.Response(writer, fmt.Sprintf("%d", id), true)
}

// FindAllProcdefPage find by page
// 分页查询
func FindAllProcdefPage(writer http.ResponseWriter, request *http.Request) {
	var procdef = service.Procdef{PageIndex: 1, PageSize: 10}
	err := util.Body2Struct(request, &procdef)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	datas, err := procdef.FindAllPageAsJSON()
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	fmt.Fprintf(writer, "%s", datas)
}

// DelProcdefByID del by id
// 根据 id 删除
func DelProcdefByID(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	var ids []string = request.Form["id"]
	if len(ids) == 0 {
		util.ResponseErr(writer, "request param 【id】 is not valid , id 不存在 ")
		return
	}
	id, err := strconv.Atoi(ids[0])
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	err = service.DelProcdefByID(id)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	util.ResponseOk(writer)
}
