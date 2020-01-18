package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	abra "github.com/ace-teknologi/abra/abra-lib"
)

var (
	// @TODO replace with interface
	abrClient *abra.Client
	abrGUID   string
)

func init() {
	abrGUID = os.Getenv("ABRA_GUID")
	var err error
	abrClient, err = abra.NewWithGuid(abrGUID)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	addr := ":" + os.Getenv("PORT")
	log.Printf("[INFO] Listening on %s", addr)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DEBUG] received request %+v", r)

	r.ParseForm()
	text := r.Form.Get("text")

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	io.WriteString(w, fmt.Sprintf("Abra is searching for %v\n", text))

	resp, err := abrClient.Search(text)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, fmt.Sprintf("An error occured: %v", err))
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Found!\n")
	io.WriteString(w, fmt.Sprintf("ABN: %v\n", resp.ABN()))
	io.WriteString(w, fmt.Sprintf("Name: %v\n", resp.Name()))
}
