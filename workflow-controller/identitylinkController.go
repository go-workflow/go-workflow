package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-workflow/go-workflow/workflow-engine/service"
	"github.com/mumushuiding/util"
)

// FindParticipantByProcInstID 根据流程id查询流程参与者
func FindParticipantByProcInstID(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		util.ResponseErr(writer, "只支持get方法！！")
		return
	}
	request.ParseForm()
	if len(request.Form["procInstID"]) == 0 {
		util.ResponseErr(writer, "流程 procInstID 不能为空")
		return
	}
	procInstID, err := strconv.Atoi(request.Form["procInstID"][0])
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	result, err := service.FindParticipantByProcInstID(procInstID)
	if err != nil {
		util.ResponseErr(writer, err)
		return
	}
	fmt.Fprintf(writer, result)

}
