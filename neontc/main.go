package main

import (
	"fmt"
	"github.com/codemicro/go-neon/neontc/tc"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("unhandled error: %v\n", err)
	}
}

func run() error {
	fmt.Println(os.Getenv("GOFILE"))
	fmt.Println(os.Getenv("GOLINE"))
	fmt.Println(os.Getenv("GOPACKAGE"))
	fmt.Println(os.Getenv("GOROOT"))

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := tc.RunOnDirectory(cwd); err != nil {
		return err
	}

	return nil
}
