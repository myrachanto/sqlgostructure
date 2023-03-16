package main

import (
	"log"

	"github.com/myrachanto/sqlgostructure/src/routes"
)

func init() {
	log.SetPrefix("SQlGoStructure :--")
}
func main() {
	log.Println("server started..........")
	routes.ApiServer()

}
