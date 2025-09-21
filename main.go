package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jangidRkt08/mongoapi/router"
)

// mongodb+srv://<db_username>:<db_password>@cluster0.x2js4d6.mongodb.net/

func main(){
	fmt.Println("MongoDb API with golang")
	fmt.Println("Server is running....")
	r:= router.Router()


	log.Fatal(http.ListenAndServe(":8080", r))
	fmt.Println("Listening at port 8080")
}
