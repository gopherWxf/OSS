package temp

/*
	POST
	$STORAGE_ROOT/temp/t.Uuid ,保存临时对象信息(uuid,name,size)
	$STORAGE_ROOT/temp/t.Uuid.dat ,用于保存临时对象内容(例如"this is test4 version 1")
*/
import (
	"OSS/lib/golog"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//记录临时对象的uuid，名字和大小
type tempInfo struct {
	Uuid string
	Name string //hashAndId
	Size int64
}

func Post(ctx *gin.Context) {
	r := ctx.Request
	w := ctx.Writer
	/*
		//windows下没有uuidgen的命令，这里使用google/uuid
		output, err := exec.Command("uuidgen").Output()
		if err!=nil{
			fmt.Println("uuid err",err)
		}
	*/
	//TODO 等到项目最终部署到linux时替换uuidgen
	//生成一个随机的uuid
	uuid, err := uuid.NewUUID()
	if err != nil {
		fmt.Println("uuid err", err)
	}
	//获取hash值
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	//获取大小
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := tempInfo{
		Uuid: uuid.String(),
		Name: name,
		Size: size,
	}
	//将结构体的内容写入磁盘 $STORAGE_ROOT/temp/t.Uuid ,保存临时对象信息
	err = t.writeToFile()
	if err != nil {
		golog.Error.Println("write to file", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//创建$STORAGE_ROOT/temp/t.Uuid.dat ，用于保存临时对象内容
	datafile, err := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid + ".dat")
	if err != nil {
		golog.Error.Println("create file err: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer datafile.Close()
	//返回临时对象的uuid，意味着数据服务节点已经为这个临时对象做好准备了
	//等待PATCH方法将数据上传
	w.Write([]byte(uuid.String()))
}

//将结构体的内容写入磁盘 $STORAGE_ROOT/temp/t.Uuid ,保存临时对象信息
func (t *tempInfo) writeToFile() error {
	file, err := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid)
	if err != nil {
		golog.Error.Println("write to file err", err, t.Uuid)
		return err
	}
	defer file.Close()
	b, _ := json.Marshal(t)
	file.Write(b)
	return nil
}
