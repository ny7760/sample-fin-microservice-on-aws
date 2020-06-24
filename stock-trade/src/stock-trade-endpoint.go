package main

import (
	"os"

	"./subfunc"

	"github.com/gin-gonic/gin"
)

// BasePass ...
var BasePass string = "trade"

// WrapperFunc ...
func WrapperFunc(fn func(c *gin.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		fn(c)
	}
}

func main() {
	PortNum := ":" + os.Getenv("PORT_NUMBER")
	r := gin.Default()

	// endpoints
	r.POST(BasePass+"/order", WrapperFunc(subfunc.OrderTrades))

	r.Run(PortNum)
}
