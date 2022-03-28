package objectStream

//创建objects的输入流和输出流
import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

/*
	PutStream在这一版本已经不用了，全部换为TempPutStream
*/
//输出流
type PutStream struct {
	writer *io.PipeWriter
	ch     chan error
}

//向流中写入数据
func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

//关闭流，查看是否存在错误信息
func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.ch
}

//返回一个writer流,写入writer的内容可以从reader中读出来，而reader作为NewRequest的参数传递进去
func NewPutStream(serverAddr, objectName string) *PutStream {
	//创建一个pipe管道
	reader, writer := io.Pipe()
	//创建一个传递错误err的管道
	ch := make(chan error)

	go func() {
		//创建一个put的request
		request, _ := http.NewRequest("PUT", "http://"+serverAddr+"/objects/"+objectName, reader)
		//创建一个客户端
		client := http.Client{}
		//客户端发送这个request。
		//这里的reader是阻塞的，所以这个Do是阻塞的
		//因为pipe中无数据，等待writer发送完消息后，reader开始工作，pipe是读写阻塞的
		//这里的reader读出来的内容就是request.body中的内容
		res, err := client.Do(request)
		defer res.Body.Close()
		//如果数据节点返回的不是200ok，则将这个错误上报
		if err == nil && res.StatusCode != http.StatusOK {
			err = fmt.Errorf("dataServer return http code %d", res.StatusCode)
		}
		ch <- err
	}()

	return &PutStream{
		writer: writer,
		ch:     ch,
	}
}

//读取流
type GetStream struct {
	io.Reader
}

//从流中读出数据
func (r *GetStream) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

//创建一个读取流，这里设计newGetStream的目的是通过检查两个参数，隐藏内部细节，不导出
func NewGetStream(serverAddr, hashAndId string) (*GetStream, error) {
	if serverAddr == "" || hashAndId == "" {
		return nil, fmt.Errorf("invalid server %s hashAndId %s", serverAddr, hashAndId)
	}
	return newGetStream("http://" + serverAddr + "/objects/" + hashAndId)
}
func newGetStream(url string) (*GetStream, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer GET return http code %d", res.StatusCode)
	}
	return &GetStream{res.Body}, nil
}

//临时对象
type TempPutStream struct {
	Server string
	Uuid   string
}

//创建临时对象,对于数据服务节点来说，hash值是对象名，这样可以实现去重
func NewTempPutStream(server, hashAndId string, size int64) (*TempPutStream, error) {
	request, err := http.NewRequest("POST", "http://"+server+"/temp/"+hashAndId, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("size", strconv.Itoa(int(size)))
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	uuid, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &TempPutStream{
		Server: server,
		Uuid:   string(uuid),
	}, nil
}

//以PATCH的方法访问数据服务的temp接口，将需要写入的数据上传
func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, err := http.NewRequest("PATCH", "http://"+w.Server+"/temp/"+w.Uuid, strings.NewReader(string(p)))
	if err != nil {
		return 0, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", response.StatusCode)
	}
	return len(p), nil
}

//根据输入的参数决定用PUT or DELETE的方法访问数据服务的temp接口
func (w *TempPutStream) Commit(good bool) {
	method := "DELETE"
	if good {
		method = "PUT"
	}
	request, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return
	}
	res.Body.Close()
}
