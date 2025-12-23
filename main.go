package main

import (
	"os"

	"github.com/consol-monitoring/check_prometheus/internal/checker"
)

func main() {

	checker.Check(os.Args)
}
