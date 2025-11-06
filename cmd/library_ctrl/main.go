package main

import (
	"context"
	"log"
	"os"

	"github.com/joematpal/library_rpc/internal/library_ctrl"
)

func main() {
	app := library_ctrl.NewApp()
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
