package main

import (
	"log"

	"github.com/Dankko0w0/gospike/cli"
)

func main() {

	// 执行 CLI 命令
	if err := cli.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
