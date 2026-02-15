package main

import (
	"log"

	"elder-care-volunteer/config"
	"elder-care-volunteer/routes"
	"elder-care-volunteer/tasks"
)

func main() {
	db := config.InitDB()
	tasks.StartNoReplyChecker(db)

	r := routes.SetupRouter(db)

	log.Println("server started at :8081")
	r.Run(":8081")
}
