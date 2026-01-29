package main

import (
	"log"

	"elder-care-volunteer/config"
	"elder-care-volunteer/routes"
)

func main() {
	db := config.InitDB()
	r := routes.SetupRouter(db)

	log.Println("server started at :8080")
	r.Run(":8080")
}
