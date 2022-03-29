package core

/*
	http://dataServerIP/objects/<xxx>
	这种形式的url相应过来
	如果是GET：
		查找objects目录下所有以<hash>.<X>开头的文件
		files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
		然后将文件中的内容写入响应体
		io.Copy(w, file)
*/
import "net/http"

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	//去除了处理PUT方法的put函数，因为现在的数据服务的对象上传完全依靠
	//temp接口的临时对象转正整，不再需要objects接口的put方法
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
