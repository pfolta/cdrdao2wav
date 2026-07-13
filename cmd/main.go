package main

import (
	"fmt"
	"runtime"

	"github.com/pfolta/cdrdao2wav/internal"
)

func printHeader() {
	fmt.Printf(
		"%s version %s-%s-%s (%s) - %s\n",
		"cdrdao2wav",
		internal.Version,
		runtime.GOARCH,
		runtime.GOOS,
		internal.BuildDate,
		"Copyright (C) Peter Folta",
	)
}

func main() {
	printHeader()
}
