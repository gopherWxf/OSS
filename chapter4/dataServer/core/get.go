package core

import (
	"ceph/chapter4/dataServer/locate"
	"ceph/chapter4/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	file := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(w, file)
}
func getFile(hash string) string {
	file := os.Getenv("STORAGE_ROOT") + "/objects/" + hash
	f, _ := os.Open(file)
	defer f.Close()
	d := url.PathEscape(utils.CalculateHash(f))
	if d != hash {
		log.Println("object hash mismatch,remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	return file
}
func sendFile(w io.Writer, file string) {
	f, _ := os.Open(file)
	defer f.Close()
	io.Copy(w, f)
}
