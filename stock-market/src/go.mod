module fin-micro/stock-market

go 1.14

replace (
	fin-micro/stock-market/subfunc => ./subfunc
	fin-micro/stock-market/subfunc/common => ./subfunc/common

)

require (
	fin-micro/stock-market/subfunc v0.0.0-00010101000000-000000000000
	fin-micro/stock-market/subfunc/common v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.32.3
	github.com/aws/aws-xray-sdk-go v1.1.0 // indirect
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/oroshnivskyy/go-gin-aws-x-ray v0.1.1 // indirect
)
