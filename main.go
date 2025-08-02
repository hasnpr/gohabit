package main

import (
	"log"

	"github.com/hasnpr/gohabit/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("can't execute app. error: %v\n", err)
	}
}
