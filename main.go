package main

import (
	"runtime"

	"github.com/charles-d-burton/goshot/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
