package kafka

import (
	"github.com/Shopify/sarama"
	"fmt"
	"logagent/log_transfer/es"
)

//初始化kafka消费者， 从kafka取数据发往ES
func Init(addrs []string, topic string) error {
	consumer, err := sarama.NewConsumer(addrs, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err:%v \n", err)
		return err
	}
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Printf("fail to get list of partition,err:%v \n", err)
		return err
	}
	fmt.Println("分区列表",partitionList)
	for partition := range partitionList {
		//针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start cnsumer for partition,err:%v \n", err)
			return err
		}

		//异步从每个分区消费信息
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages(){
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
				ld := es.LogData{
					Topic:topic,
					Data:string(msg.Value),
				}
				//es.SendToES(topic, ld) // 函数调用函数
				//优化： 把数据封装并传到chan中
				es.SendTOCHan(&ld)
			}
		}(pc)

	}
	return err
}