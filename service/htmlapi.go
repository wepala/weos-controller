package service

import (
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
)

func NewHTMLServer(service ServiceInterface, staticFolder string) http.Handler {
	//TODO setup a handler using gorilla + negroni (or just negroni?)
	router := mux.NewRouter()
	router.PathPrefix("/static").Handler(negroni.New(negroni.NewStatic(http.Dir(staticFolder))))
	n := negroni.Classic()
	//TODO add middleware that returns the html response
	n.UseHandler(router)
	return n
}
