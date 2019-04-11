/**
 * main process
 * Usage: wx-server[ -v]
 * Rosbit Xu
 */
package main

import (
	"os"
	"fmt"
)

// variables set via go build -ldflags
var (
	buildTime string
	osInfo    string
	goInfo    string
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		ShowInfo("name",       os.Args[0])
		ShowInfo("build time", buildTime)
		ShowInfo("os name",    osInfo)
		ShowInfo("compiler",   goInfo)
		return
	}

	err := CheckGlobalConf()
	if err != nil {
		fmt.Printf("Failed to check conf: %v\n", err)
		os.Exit(3)
		return
	}
	DumpConf()

	StartWxApp()
	os.Exit(0)
}

func ShowInfo(prompt, info string) {
	if info != "" {
		fmt.Printf("%10s: %s\n", prompt, info)
	}
}
