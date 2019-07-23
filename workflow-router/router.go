package router

import (
	"net/http"

	config "github.com/go-workflow/go-workflow/workflow-config"
	controller "github.com/go-workflow/go-workflow/workflow-controller"
)

// Mux 路由
var Mux = http.NewServeMux()
var conf = *config.Config

func init() {
	setMux()
}
func intercept(h http.HandlerFunc) http.HandlerFunc {
	return crossOrigin(h)
}
func crossOrigin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", conf.AccessControlAllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", conf.AccessControlAllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", conf.AccessControlAllowHeaders)
		h(w, r)
	}
}
func setMux() {
	Mux.HandleFunc("/api/v1/workflow/", controller.Index)
	//-------------------------流程定义----------------------
	Mux.HandleFunc("/api/v1/workflow/procdef/save", intercept(controller.SaveProcdef))
	Mux.HandleFunc("/api/v1/workflow/procdef/saveByToken", intercept(controller.SaveProcdefByToken))
	Mux.HandleFunc("/api/v1/workflow/procdef/findAll", intercept(controller.FindAllProcdefPage))
	Mux.HandleFunc("/api/v1/workflow/procdef/delById", intercept(controller.DelProcdefByID))
	// -----------------------流程实例-----------------------
	Mux.HandleFunc("/api/v1/workflow/process/start", intercept(controller.StartProcessInstance))               // 启动流程
	Mux.HandleFunc("/api/v1/workflow/process/startByToken", intercept(controller.StartProcessInstanceByToken)) // 启动流程
	Mux.HandleFunc("/api/v1/workflow/process/findTask", intercept(controller.FindMyProcInstPageAsJSON))        // 查询需要我审批的流程
	Mux.HandleFunc("/api/v1/workflow/process/findTaskByToken", intercept(controller.FindMyProcInstByToken))
	Mux.HandleFunc("/api/v1/workflow/process/startByMyself", intercept(controller.StartByMyself))   // 查询我启动的流程
	Mux.HandleFunc("/api/v1/workflow/process/FindProcNotify", intercept(controller.FindProcNotify)) // 查询抄送我的流程
	// Mux.HandleFunc("/workflow/process/moveToHistory", controller.MoveFinishedProcInstToHistory)
	// -----------------------任务--------------------------
	Mux.HandleFunc("/api/v1/workflow/task/complete", intercept(controller.CompleteTask))
	Mux.HandleFunc("/api/v1/workflow/task/completeByToken", intercept(controller.CompleteTaskByToken))
	Mux.HandleFunc("/api/v1/workflow/task/withdraw", intercept(controller.WithDrawTask))
	Mux.HandleFunc("/api/v1/workflow/task/withdrawByToken", intercept(controller.WithDrawTaskByToken))
	// ----------------------- 关系表 -------------------------
	Mux.HandleFunc("/api/v1/workflow/identitylink/findParticipant", intercept(controller.FindParticipantByProcInstID))

	// ******************************** 历史纪录 ***********************************
	// -------------------------- 流程实例 -------------------------------
	Mux.HandleFunc("/api/v1/workflow/procHistory/findTask", intercept(controller.FindProcHistory))
	Mux.HandleFunc("/api/v1/workflow/procHistory/findTaskByToken", intercept(controller.FindProcHistoryByToken))
	Mux.HandleFunc("/api/v1/workflow/procHistory/startByMyself", intercept(controller.StartHistoryByMyself))   // 查询我启动的流程
	Mux.HandleFunc("/api/v1/workflow/procHistory/FindProcNotify", intercept(controller.FindProcHistoryNotify)) // 查询抄送我的流程
	// ----------------------- 关系表 -------------------------
	Mux.HandleFunc("/api/v1/workflow/identitylinkHistory/findParticipant", intercept(controller.FindParticipantHistoryByProcInstID))

	// ************************** 以下为废弃接口 *****************************************************
	//**********************************************************************************************
	Mux.HandleFunc("/workflow/", controller.Index)
	//-------------------------流程定义----------------------
	Mux.HandleFunc("/workflow/procdef/save", intercept(controller.SaveProcdef))
	Mux.HandleFunc("/workflow/procdef/saveByToken", intercept(controller.SaveProcdefByToken))
	Mux.HandleFunc("/workflow/procdef/findAll", intercept(controller.FindAllProcdefPage))
	Mux.HandleFunc("/workflow/procdef/delById", intercept(controller.DelProcdefByID))
	// -----------------------流程实例-----------------------
	Mux.HandleFunc("/workflow/process/start", intercept(controller.StartProcessInstance))               // 启动流程
	Mux.HandleFunc("/workflow/process/startByToken", intercept(controller.StartProcessInstanceByToken)) // 启动流程
	Mux.HandleFunc("/workflow/process/findTask", intercept(controller.FindMyProcInstPageAsJSON))        // 查询需要我审批的流程
	Mux.HandleFunc("/workflow/process/findTaskByToken", intercept(controller.FindMyProcInstByToken))
	Mux.HandleFunc("/workflow/process/startByMyself", intercept(controller.StartByMyself))   // 查询我启动的流程
	Mux.HandleFunc("/workflow/process/FindProcNotify", intercept(controller.FindProcNotify)) // 查询抄送我的流程
	// Mux.HandleFunc("/workflow/process/moveToHistory", controller.MoveFinishedProcInstToHistory)
	// -----------------------任务--------------------------
	Mux.HandleFunc("/workflow/task/complete", intercept(controller.CompleteTask))
	Mux.HandleFunc("/workflow/task/completeByToken", intercept(controller.CompleteTaskByToken))
	Mux.HandleFunc("/workflow/task/withdraw", intercept(controller.WithDrawTask))
	Mux.HandleFunc("/workflow/task/withdrawByToken", intercept(controller.WithDrawTaskByToken))
	// ----------------------- 关系表 -------------------------
	Mux.HandleFunc("/workflow/identitylink/findParticipant", intercept(controller.FindParticipantByProcInstID))

	// ******************************** 历史纪录 ***********************************
	// -------------------------- 流程实例 -------------------------------
	Mux.HandleFunc("/workflow/procHistory/findTask", intercept(controller.FindProcHistory))
	Mux.HandleFunc("/workflow/procHistory/findTaskByToken", intercept(controller.FindProcHistoryByToken))
	Mux.HandleFunc("/workflow/procHistory/startByMyself", intercept(controller.StartHistoryByMyself))   // 查询我启动的流程
	Mux.HandleFunc("/workflow/procHistory/FindProcNotify", intercept(controller.FindProcHistoryNotify)) // 查询抄送我的流程
	// ----------------------- 关系表 -------------------------
	Mux.HandleFunc("/workflow/identitylinkHistory/findParticipant", intercept(controller.FindParticipantHistoryByProcInstID))
}
