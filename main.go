package main

import (
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/api"
)

func main() {
	server := api.NewApiServer(http.NewServeMux())
	err := http.ListenAndServe(":8080", server)
	if err != nil {
		panic(err.Error())
	}
}
