module github.com/consol-monitoring/check_prometheus

go 1.25

replace internal/helper => ./internal/helper

replace internal/mode => ./internal/mode

require (
	github.com/consol-monitoring/check_x v0.0.0-20230423195421-be7cfdc8c478
	github.com/urfave/cli v1.22.14
	internal/helper v0.0.0-00010101000000-000000000000
	internal/mode v0.0.0-00010101000000-000000000000
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/prometheus/client_golang v1.18.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
)
