package main

import (
	"os"

	"./subfunc"

	"github.com/gin-gonic/gin"
)

// BasePass ...
var BasePass string = "stock"

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
	r.GET(BasePass+"/attribute", WrapperFunc(subfunc.GetAttribute))
	r.GET(BasePass+"/prices", WrapperFunc(subfunc.GetHistoricalPrices))

	r.Run(PortNum)
}
