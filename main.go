package main

import (
	"os"

	"github.com/consol-monitoring/check_prometheus/pkg/checker"
)

func main() {
	return_code := checker.CheckMain(os.Args)

	os.Exit(return_code)
}
