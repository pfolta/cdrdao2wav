package main

import (
	"fmt"
	"runtime"

	"github.com/pfolta/cdrdao2audio/internal"
)

func printHeader() {
	fmt.Printf(
		"%s version %s-%s-%s (%s) - %s\n",
		"cdrdao2audio",
		internal.Version,
		runtime.GOOS,
		runtime.GOARCH,
		internal.BuildDate,
		"Copyright (C) Peter Folta",
	)
}

func main() {
	printHeader()
}
