package main

import (
	"os"

	"github.com/consol-monitoring/check_prometheus/pkg/checker"
)

func main() {

	checker.Check(os.Args)
}
