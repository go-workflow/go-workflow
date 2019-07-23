package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	router "github.com/go-workflow/go-workflow/workflow-router"

	config "github.com/go-workflow/go-workflow/workflow-config"

	model "github.com/go-workflow/go-workflow/workflow-engine/model"
	"github.com/go-workflow/go-workflow/workflow-engine/service"
)

// 配置
var conf = *config.Config

func crossOrigin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", conf.AccessControlAllowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", conf.AccessControlAllowMethods)
		w.Header().Set("Access-Control-Allow-Headers", conf.AccessControlAllowHeaders)
		h(w, r)
	}
}
func main() {
	mux := router.Mux
	// 启动数据库连接
	model.Setup()
	// 启动redis连接
	model.SetRedis()
	// 启动定时任务
	service.CronJobs()
	// 启动服务
	readTimeout, err := strconv.Atoi(conf.ReadTimeout)
	if err != nil {
		panic(err)
	}
	writeTimeout, err := strconv.Atoi(conf.WriteTimeout)
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", conf.Port),
		Handler:        mux,
		ReadTimeout:    time.Duration(readTimeout * int(time.Second)),
		WriteTimeout:   time.Duration(writeTimeout * int(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("the application start up at port%s", server.Addr)
	if conf.TLSOpen == "true" {
		err = server.ListenAndServeTLS(conf.TLSCrt, conf.TLSKey)
	} else {
		err = server.ListenAndServe()
	}
	if err != nil {
		log.Printf("Server err: %v", err)
	}

}
