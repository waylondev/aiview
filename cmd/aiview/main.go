package main

import (
	"os"

	"github.com/jackwener/aiview/internal/cli"
	_ "github.com/jackwener/aiview/internal/platform/bilibili"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}