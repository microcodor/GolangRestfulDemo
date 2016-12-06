package main

import (
	"log"
	"net/http"
)

/*
microcodor app 接口server
createed 2016/12/05
*/
func main() {

	router := NewRouter()

	log.Fatal(http.ListenAndServe(":9090", router))
}
