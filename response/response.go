package response

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func PrintResponse(w http.ResponseWriter, val string, err []string) []byte {
	return []byte(fmt.Sprintf("%s %s", val, strings.Join(err[:], ",")))
}

func InternalServerError(w http.ResponseWriter, err ...string) {
	w.WriteHeader(http.StatusInternalServerError)
	PrintResponse(w, "500 - Server Error", err)
	log.Println("Error", err)
}

func NotFoundError(w http.ResponseWriter, err ...string) {
	w.WriteHeader(http.StatusNotFound)
	PrintResponse(w, "404 - Not Found", err)
}

func BadRequestError(w http.ResponseWriter, err ...string) {
	w.WriteHeader(http.StatusBadRequest)
	PrintResponse(w, "400 - Bad Request", err)
}
