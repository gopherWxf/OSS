package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type Metadata struct {
	Name    string
	Version string
	Size    string
	Hash    string
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	request, err := http.Get("http://" + "10.29.2.1:12345" + "/versions/")
	if err != nil {
		log.Println(err)
		return
	}
	s := bufio.NewScanner(request.Body)

	w.Write([]byte(`<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
    <style>
        body{
            background-color: rgb(237, 242, 245);
        }
        div{
            width: 500px;
            height: 500px;
            margin: 0 auto;
            text-align: center;
        }
        table{
            margin: 0 auto;
        }
        a{
            text-decoration: none;
            color: #009ACD;
        }
        a:hover{
            color: green;
            font-size: 26px;
        }
    </style>
</head>

<body >
    <div>
        <table>
            <tr>
                <th>文件名</th>
                <th>版本</th>
                <th>大小</th>
            </tr>`))
	for s.Scan() {
		var meta Metadata
		json.Unmarshal([]byte(s.Text()), &meta)
		if meta.Hash != "" {
			//中文
			n, _ := url.PathUnescape(meta.Name)
			size, _ := strconv.ParseFloat(meta.Size, 0)
			size /= 1024
			l := fmt.Sprintf("<tr><td><a href=/download?name=%s&version=%s>%s</a></td><td>%s</td><td>%.3fkb</td></tr>", url.PathEscape(n), meta.Version, n, meta.Version, size)
			w.Write([]byte(l))
		}
	}
	w.Write([]byte(`</table>
        <form action=/upload method=post enctype=multipart/form-data>
            <input type=file name=upload>
            <input type=submit>
        </form>
    </div>
</body>

</html>`))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	f, header, e := r.FormFile("upload")
	if e != nil {
		log.Println(e)
		return
	}
	defer f.Close()
	h := sha256.New()
	io.Copy(h, f)
	d := base64.StdEncoding.EncodeToString(h.Sum(nil))
	log.Println(d)
	f.Seek(0, 0)
	dat, _ := ioutil.ReadAll(f)
	req, e := http.NewRequest("PUT", "http://"+"10.29.2.1:12345"+"/objects/"+url.PathEscape(header.Filename), bytes.NewBuffer(dat))
	if e != nil {
		log.Println(e)
		return
	}
	req.Header.Set("digest", "SHA-256="+d)
	client := http.Client{}
	log.Println("uploading file", header.Filename, "hash", d, "size", header.Size)
	_, e = client.Do(req)
	if e != nil {
		log.Println(e)
		return
	}
	log.Println("uploaded")
	time.Sleep(time.Second)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	req, e := http.Get("http://" + "10.29.2.1:12345" + "/objects/" + url.PathEscape(r.URL.Query()["name"][0]) + "?version=" + r.URL.Query()["version"][0])
	if e != nil {
		log.Println(e)
		return
	}
	w.Header().Set("content-disposition", "attachment;filename="+r.URL.Query()["name"][0])
	io.Copy(w, req.Body)
}
