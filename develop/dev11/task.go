package main

import (
	"context"
	"dev11/app"
)

func main() {
	calendar := app.NewApp()

	calendar.Run(context.Background())
}
