package core

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	file, err := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.String(), "/")[2])
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	io.Copy(file, r.Body)
	log.Println(os.Getenv("LISTEN_ADDRESS"), os.Getenv("STORAGE_ROOT")+"/objects/"+strings.Split(r.URL.String(), "/")[2], "Stored")
}
