package framework

import (
	"fmt"
	"net/http"
)

type Engine struct{

}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.RequestURI)
}