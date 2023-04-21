module github.com/consol-monitoring/check_prometheus

go 1.19

replace internal/helper => ./internal/helper

replace internal/mode => ./internal/mode

require (
	internal/helper v0.0.0-00010101000000-000000000000
	internal/mode v0.0.0-00010101000000-000000000000
)
