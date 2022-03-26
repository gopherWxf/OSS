package objectStream

//创建objects的输入流和输出流
import (
	"fmt"
	"io"
	"net/http"
)

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
func NewGetStream(serverAddr, objectName string) (*GetStream, error) {
	if serverAddr == "" || objectName == "" {
		return nil, fmt.Errorf("invalid server %s object %s", serverAddr, objectName)
	}
	return newGetStream("http://" + serverAddr + "/objects/" + objectName)
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
