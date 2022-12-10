package main

import (
	"JOLFramework/framework"
	"fmt"
	"log"
	"net/http"
)

func main() {
	engine := framework.Engine{}
	port := 8080
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d",port), &engine))
	fmt.Println("listening on port:", port)
}
