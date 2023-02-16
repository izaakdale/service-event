package main

import (
	"github.com/izaakdale/service-event/internal/app"
	_ "github.com/lib/pq"
)

func main() {
	app.Run()
}
