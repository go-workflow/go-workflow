package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	config "github.com/go-workflow/go-workflow/workflow-config"
	controller "github.com/go-workflow/go-workflow/workflow-controller"
	model "github.com/go-workflow/go-workflow/workflow-engine/model"
	"github.com/go-workflow/go-workflow/workflow-engine/service"
)

func crossOrigin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h(w, r)
	}
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/workflow/", controller.Index)
	//-------------------------流程定义----------------------
	mux.HandleFunc("/workflow/procdef/save", controller.SaveProcdef)
	mux.HandleFunc("/workflow/procdef/findAll", controller.FindAllProcdefPage)
	mux.HandleFunc("/workflow/procdef/delById", controller.DelProcdefByID)
	// -----------------------流程实例-----------------------
	mux.HandleFunc("/workflow/process/start", controller.StartProcessInstance)               // 启动流程
	mux.HandleFunc("/workflow/process/startByToken", controller.StartProcessInstanceByToken) // 启动流程
	mux.HandleFunc("/workflow/process/findTask", controller.FindMyProcInstPageAsJSON)        // 查询需要我审批的流程
	mux.HandleFunc("/workflow/process/findTaskByToken", controller.FindMyProcInstByToken)
	mux.HandleFunc("/workflow/process/startByMyself", controller.StartByMyself)   // 查询我启动的流程
	mux.HandleFunc("/workflow/process/FindProcNotify", controller.FindProcNotify) // 查询抄送我的流程
	// mux.HandleFunc("/workflow/process/moveToHistory", controller.MoveFinishedProcInstToHistory)
	// -----------------------任务--------------------------
	mux.HandleFunc("/workflow/task/complete", controller.CompleteTask)
	mux.HandleFunc("/workflow/task/completeByToken", controller.CompleteTaskByToken)
	mux.HandleFunc("/workflow/task/withdraw", controller.WithDrawTask)
	mux.HandleFunc("/workflow/task/withdrawByToken", controller.WithDrawTaskByToken)
	// ----------------------- 关系表 -------------------------
	mux.HandleFunc("/workflow/identitylink/findParticipant", controller.FindParticipantByProcInstID)

	// ******************************** 历史纪录 ***********************************
	// -------------------------- 流程实例 -------------------------------
	mux.HandleFunc("/workflow/procHistory/findTask", controller.FindProcHistory)
	mux.HandleFunc("/workflow/procHistory/findTaskByToken", controller.FindProcHistoryByToken)
	mux.HandleFunc("/workflow/procHistory/startByMyself", controller.StartHistoryByMyself)   // 查询我启动的流程
	mux.HandleFunc("/workflow/procHistory/FindProcNotify", controller.FindProcHistoryNotify) // 查询抄送我的流程
	// ----------------------- 关系表 -------------------------
	mux.HandleFunc("/workflow/identitylinkHistory/findParticipant", controller.FindParticipantHistoryByProcInstID)
	// 配置
	var config = *config.Config
	// 启动数据库连接
	model.Setup()
	// 启动redis连接
	model.SetRedis()
	// 启动定时任务
	service.CronJobs()
	// 启动服务
	readTimeout, err := strconv.Atoi(config.ReadTimeout)
	if err != nil {
		panic(err)
	}
	writeTimeout, err := strconv.Atoi(config.WriteTimeout)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(readTimeout * int(time.Second)),
		WriteTimeout:   time.Duration(writeTimeout * int(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("the application start up at port%s", server.Addr)
	if config.TLSOpen == "true" {
		err = server.ListenAndServeTLS(config.TLSCrt, config.TLSKey)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Printf("Server err: %v", err)
	}

}
