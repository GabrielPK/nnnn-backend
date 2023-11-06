package handler

import (
	// "encoding/json"
	"net/http"
	"fmt"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}