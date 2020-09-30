package kafka

import (
	"github.com/Shopify/sarama"
	"fmt"
)

var (
	client sarama.SyncProducer
)

//基于sarama第三方库的kafka client
func Init(addrs []string) (err error) {
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
	fmt.Println("连接kafka成功")
	return
}

func SendToKafka(topic, data string)  {
	//构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	//发送消息到kafka
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed err:", err)
		return
	}
	fmt.Printf("pid:%v, offset:%v\n",pid, offset)
	fmt.Println("发送成功")
}