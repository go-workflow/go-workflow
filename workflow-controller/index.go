package controller

import (
	"fmt"
	"net/http"
)

// Index 首页
func Index(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(writer, "Hello world!")
}
