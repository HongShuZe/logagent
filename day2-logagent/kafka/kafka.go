package kafka

import (
	"github.com/Shopify/sarama"
	"fmt"
	"time"
)

type logData struct {
	topic string
	data string
}

var (
	client sarama.SyncProducer // 声明一个全局的连接kafka的生产者client
	logDataChan chan *logData
)

//基于sarama第三方库的kafka client
func Init(addrs []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	//tail包的使用
	config.Producer.RequiredAcks = sarama.WaitForAll //发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner //新选出一个partition
	config.Producer.Return.Successes = true //成功交付的消息将在success channel返回


	//连接kafka
	client, err = sarama.NewSyncProducer(addrs, config)
	if err != nil {
		fmt.Println("producer close err :", err)
		return
	}
	//初始化lodDataChan
	logDataChan = make(chan *logData, maxSize)
	//开启后台的goroutine从通道中取取数据发往kafka
	//fmt.Println("连接kafka成功")
	go sendToKafka()
	return
}

//给外部暴露的一个函数，该函数只把日志数据发送到一个内部的channel中
func SendToKafka(topic, data string)  {
	msg := &logData{
		topic: topic,
		data: data,
	}
	logDataChan <- msg
}

// 真正往kafka发送日志的函数
func sendToKafka() {
	for  {
		select {
			case ld := <-logDataChan:
				//构造一个消息
				msg := &sarama.ProducerMessage{}
				msg.Topic = ld.topic
				msg.Value = sarama.StringEncoder(ld.data)

				//发送消息到kafka
				pid, offset, err := client.SendMessage(msg)
				if err != nil {
					fmt.Println("send msg failed err:", err)
					return
				}
				fmt.Printf("pid:%v, offset:%v\n",pid, offset)
				fmt.Println("发送成功")
		default:
			time.Sleep(time.Microsecond*50)
		}
	}
}