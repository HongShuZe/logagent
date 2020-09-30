package main

import (
	"github.com/Shopify/sarama"
	"fmt"
)

func main()  {
	consumer, err := sarama.NewConsumer([]string{"127.0.0.1:9092"}, nil)
	if err != nil {
		fmt.Printf("fail to start consumer, err %v\n", err)
		return
	}

	partitionList, err := consumer.Partitions("web_log")
	if err != nil {
		fmt.Printf("fail to start partition, err %v\n", err)
	}

	fmt.Println("分区列表", partitionList)
	for partition := range partitionList{
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("fail to start consumer for partition, err %v\n", err)
			return
		}
		defer pc.AsyncClose()
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			}
		}(pc)
	}
	select {
	}

}
