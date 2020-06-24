package main

import (
	"net/http"
	"os"

	"./subfunc"

	"github.com/aws/aws-xray-sdk-go/xray"
	// "github.com/gin-gonic/gin"
)

var (
	ServiceName = os.Getenv("SERVICE_NAME")
)

// BasePass ...
var BasePass string = "balance"

// WrapperFunc ... if use gin
// func WrapperFunc(fn func(c *gin.Context)) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fn(c)
// 	}
// }

func main() {
	PortNum := ":" + os.Getenv("PORT_NUMBER")

	// if use gin
	// r := gin.Default()
	// r.POST(BasePass, WrapperFunc(subfunc.UpdateBalance))
	// r.Run(PortNum)

	http.Handle("/"+BasePass, xray.Handler(xray.NewFixedSegmentNamer(ServiceName), http.HandlerFunc(subfunc.UpdateBalance)))
	http.ListenAndServe(PortNum, nil)
}
