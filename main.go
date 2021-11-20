package main

import (
	"fmt"
	"os"

	"github.com/LilRooks/five/internal/app"
)

func main() {
	if code, err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(code)
	}
}
