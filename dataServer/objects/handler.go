package objects

//func Handler(w http.ResponseWriter, r *http.Request) {
//	m := r.Method
//	//去除了处理PUT方法的put函数，因为现在的数据服务的对象上传完全依靠
//	//temp接口的临时对象转正整，不再需要objects接口的put方法
//	if m == http.MethodGet {
//		get(w, r)
//		return
//	}
//	if m == http.MethodDelete {
//		del(w, r)
//		return
//	}
//	w.WriteHeader(http.StatusMethodNotAllowed)
//}
