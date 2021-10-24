package main

import (
	"log"

	"net/http"

	config "example.com/gorestapi/configs"

	routes "example.com/gorestapi/routes"
)

func main() {
	// Connect DB
	config.ConnectDBandProperties()

	// Init Router

	// Route Handlers / Endpoints
	http.Handle("/", routes.Handlers())

	log.Fatal(http.ListenAndServe(":"+"4747", nil))
}
