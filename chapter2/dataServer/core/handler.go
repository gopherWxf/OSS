package core

/*
	http://dataServerIP/objects/<xxx>
	这种形式的url相应过来
	如果是PUT：
		首先创建一个object文件
		os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.String(), "/")[2])
		然后将请求里的body存入文件
		io.Copy(file, r.Body)
	如果是GET：
		首先打开这个文件
		os.Open(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
		然后将文件中的内容写入响应体
		io.Copy(w, file)
*/
import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
