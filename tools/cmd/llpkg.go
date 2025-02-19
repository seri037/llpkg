package main

import (
	"flag"
	"fmt"

	configCLI "github.com/goplus/llpkg/llpkg-tool/pkg/config"
)

var llpkgConfigPath = flag.String("config", "", "path to config file")

func main() {
	config, err := configCLI.ParseLLpkgConfig("./llpkg-tool/demo/.llpkg/llpkg.cfg")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(config)
	flag.Parse()
	if *llpkgConfigPath == "" {
		printUsage()
		return
	}

	println("TODO")
}

func printUsage() {
	println("Usage: llpkg -config <path to config file>")
}
