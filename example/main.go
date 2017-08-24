package main

import (
	"fmt"
	"net/http"

	nss "github.com/astj/negroni-server-status"
	"github.com/urfave/negroni"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "This is from internal")
	})

	n := negroni.New()

	n.Use(nss.NewMiddleware("/server-status"))

	n.UseHandler(mux)

	n.Run(":3000")
}
