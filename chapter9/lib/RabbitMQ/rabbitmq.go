package RabbitMQ

//rabbitmq.go是对RabbitMQ的封装

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

// rabbitMQ结构体
type RabbitMQ struct {
	channel *amqp.Channel
	conn    *amqp.Connection
	Name    string
}

//创建结构体实例
func NewRabbitMQ(mqUrl string) *RabbitMQ {
	conn, err := amqp.Dial(mqUrl)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	//申请队列
	queue, err := ch.QueueDeclare(
		"",
		//是否持久化
		false,
		//是否自动删除
		true,
		//是否具有排他性
		false,
		//是否阻塞处理
		false,
		//额外的属性
		nil,
	)
	mq := new(RabbitMQ)
	mq.channel = ch
	mq.Name = queue.Name
	mq.conn = conn
	return mq
}

//绑定队列到exchange中
func (r *RabbitMQ) Bind(exchange string) {
	err := r.channel.QueueBind(
		r.Name,
		"",       //在pub/sub模式下，这里的key要为空
		exchange, //exchange
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
}

//向一个exchange中发布消息
func (r *RabbitMQ) Publish(exchange string, body interface{}) {
	bytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	//发布到exchange
	err = r.channel.Publish(exchange, "", false, false, amqp.Publishing{
		ReplyTo: r.Name,
		Body:    bytes,
	})
	if err != nil {
		panic(err)
	}
}

//往指定的消息队列发送消息
func (r *RabbitMQ) Send(queue string, body interface{}) {
	bytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	//发布到queue
	err = r.channel.Publish("", queue, false, false, amqp.Publishing{
		ReplyTo: r.Name,
		Body:    bytes,
	})
	if err != nil {
		panic(err)
	}
}

//消费消息，返回消费消息的channel
func (r *RabbitMQ) Consume() <-chan amqp.Delivery {
	ch, err := r.channel.Consume(r.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	return ch
}

//关闭消息队列
func (r *RabbitMQ) Close() {
	r.channel.Close()
}
